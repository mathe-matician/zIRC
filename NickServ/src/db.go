package nickserv

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/phuslu/log"

	"github.com/jackc/pgx/v5"
)

var g_DB *pgx.Conn

func db_init() {
	// var tls_config *tls.Config
	enable_tls, err := strconv.ParseBool(GetEnv("IRC_DB_ENABLE_TLS", "false"))
	if err != nil {
		log.Error().Msg(err.Error())
	}

	db_user := GetEnv("IRC_DB_USER", "postgres")
	db_database := GetEnv("IRC_DB_DATABASE", "postgres")
	db_app_name := GetEnv("IRC_DB_APPLICATION_NAME", "zirc_server")
	db_host := GetEnv("IRC_DB_HOST", "zirc_db")
	db_port := GetEnv("IRC_DB_PORT", "5432")
	sslmode := GetEnv("IRC_DB_SSLMODE", "prefer")

	connStr := fmt.Sprintf("host='%s' port='%s' dbname='%s' user='%s' sslmode='%s' application_name='%s'", db_host, db_port, db_database, db_user, sslmode, db_app_name)
	if enable_tls {
		crt_path := GetEnv("IRC_DB_TLS_CERT_PATH", "")
		key_path := GetEnv("IRC_DB_TLS_KEY_PATH", "")
		// tls_config := CreateTLSConfig(crt_path, key_path)
		connStr += fmt.Sprintf("sslcert='%s' sslkey='%s' sslrootcert='%s'", crt_path, key_path)
	}
	pgConfig, err := pgx.ParseConfig(connStr)
	if err != nil {
		log.Error().Msgf("Unable to parse db config: %s\n", err.Error())
		panic(err)
	}

	g_DB, err := pgx.ConnectConfig(context.Background(), pgConfig)
	if err != nil {
		log.Error().Msgf("Unable to connect to database: %s\n", err.Error())
		panic(err)
	}

	ctx := context.Background()
	if err := g_DB.Ping(ctx); err != nil {
		log.Info().Msgf("Could not ping database upon startup: %s", err)
	} else {
		var row string
		err := g_DB.QueryRow(context.Background(), "SELECT version()").Scan(&row)
		if err != nil {
			log.Error().Msgf("Unable to get version from db: %s\n", err.Error())
		} else {
			log.Info().Msgf("Pg version: %s", row)
		}
		log.Info().Msgf("Successfully connected to db. Host: %s, User: %s, Db: %s, AppName: %s", db_host, db_user, db_database, db_app_name)
	}
}

func reconnect_db_listener() {
	log.Info().Msgf("DB Reconnect Listener started")
	reconnect_wait_interval, err := strconv.Atoi(GetEnv("IRC_DB_HEALTHCHECK_INTERVAL", "3"))
	if err != nil {
		log.Error().Msgf(err.Error())
	}

	ctx := context.Background()
	for {
		if err := g_DB.Ping(ctx); err != nil {
			log.Info().Msgf("Lost connection to database: %s", err)
			log.Info().Msgf("Reconnecting to database...")
			db_init()
			return
		}
		time.Sleep(time.Duration(reconnect_wait_interval) * time.Second)
	}
}

func DB_pkg_init() {
	db_init()

	// var err error
	// will return this error:
	// ERROR:  there is no unique or exclusion constraint matching the ON CONFLICT specification
	// register_user_stmt, err = g_DB.Prepare(`INSERT INTO users VALUES (default, $1::text, $2::text) ON CONFLICT (nick, email) DO NOTHING`)
	// if err != nil {
	// 	log.Error().Msgf("Error creating register_user_stmt: %s", err.Error())
	// 	panic(err)
	// }

	go reconnect_db_listener()
}

type UserModel struct {
	username string
	password string
}
