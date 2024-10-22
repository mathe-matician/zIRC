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
var pgConfig *pgx.ConnConfig

func genPgConfig() {
	// var tls_config *tls.Config
	db_user := GetEnv("IRC_DB_USER", "postgres")
	db_database := GetEnv("IRC_DB_DATABASE", "postgres")
	db_app_name := GetEnv("IRC_DB_APPLICATION_NAME", "zirc_server")
	db_host := GetEnv("IRC_DB_HOST", "zirc_db")
	db_port := GetEnv("IRC_DB_PORT", "5432")
	sslmode := GetEnv("IRC_DB_SSLMODE", "disable")
	db_password := GetEnv("IRC_DB_PASSWORD", "password")

	connStr := fmt.Sprintf("host=%s port=%s dbname=%s password=%s user=%s sslmode=%s application_name=%s", db_host, db_port, db_database, db_password, db_user, sslmode, db_app_name)

	enable_tls, err := strconv.ParseBool(GetEnv("IRC_DB_ENABLE_TLS", "false"))
	if err != nil {
		log.Error().Msg(err.Error())
	}
	if enable_tls {
		crt_path := GetEnv("IRC_DB_TLS_CERT_PATH", "")
		key_path := GetEnv("IRC_DB_TLS_KEY_PATH", "")
		sslrootcert := GetEnv("IRC_DB_TLS_ROOT_CERT", "")
		// tls_config := CreateTLSConfig(crt_path, key_path)
		connStr += fmt.Sprintf(" sslcert='%s' sslkey='%s' sslrootcert='%s'", crt_path, key_path, sslrootcert)
	}
	pgConfig, err = pgx.ParseConfig(connStr)
	if err != nil {
		log.Error().Msgf("Unable to parse db config: %s\n", err.Error())
		panic(err)
	}

	log.Info().Msgf("Created config: Host: %s, User: %s, Db: %s, AppName: %s", db_host, db_user, db_database, db_app_name)
}

func db_init() {
	var err error
	for {
		g_DB, err = pgx.ConnectConfig(context.Background(), pgConfig)
		if err != nil {
			log.Error().Msgf("Unable to init connection to database: %s\n", err.Error())
			log.Error().Msgf("Waiting for DB to start...")
		} else {
			break
		}
		time.Sleep(time.Duration(3) * time.Second)
	}
	log.Info().Msgf("Successfully connected to DB!")

	ctx := context.Background()
	if err := g_DB.Ping(ctx); err != nil {
		log.Info().Msgf("Could not ping database upon startup: %s", err)
	} else {
		var row string
		err := g_DB.QueryRow(context.Background(), "SELECT version()").Scan(&row)
		if err != nil {
			log.Error().Msgf("Unable to get version from db: %s\n", err.Error())
		} else {
			log.Info().Msgf("%s", row)
		}
	}
}

func reconnect_db_listener() {
	var err error
	reconnect_wait_interval, err := strconv.Atoi(GetEnv("IRC_DB_HEALTHCHECK_INTERVAL", "3"))
	if err != nil {
		log.Error().Msgf("Error converting healthcheck interval to int: %s", err.Error())
		panic(err)
	}
	log.Info().Msgf("DB Reconnect Listener started. Healthcheck interval: %d", reconnect_wait_interval)

	ctx := context.Background()
	for {
		reconn := "Reconnecting to database..."
		unable := "Unable to reconnect to database: %s\n"
		if g_DB == nil {
			log.Info().Msgf(reconn)
			g_DB, err = pgx.ConnectConfig(context.Background(), pgConfig)
			if err != nil {
				log.Error().Msgf(unable, err.Error())
			}
		} else if err := g_DB.Ping(ctx); err != nil {
			log.Info().Msgf(reconn)
			g_DB.Close(ctx)
			g_DB, err = pgx.ConnectConfig(context.Background(), pgConfig)
			if err != nil {
				log.Error().Msgf(unable, err.Error())
			}
		} else {
			log.Info().Msgf("DB healthcheck running...")
		}
		time.Sleep(time.Duration(reconnect_wait_interval) * time.Second)
	}
}

func DB_pkg_init() {
	genPgConfig()
	db_init()
	go reconnect_db_listener()
}
