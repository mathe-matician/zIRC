package chat

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"time"
	"zirc/helpers"

	"github.com/go-pg/pg/v10"
	"github.com/phuslu/log"
)

var g_DB *pg.DB

// var plain_auth_stmt *pg.Stmt
// var register_stmt *pg.Stmt

func db_init() {
	var tls_config *tls.Config
	enable_tls, err := strconv.ParseBool(helpers.GetEnv("IRC_DB_ENABLE_TLS", "false"))
	if err != nil {
		log.Error().Msg(err.Error())
	}
	if enable_tls {
		crt_path := helpers.GetEnv("IRC_DB_TLS_CERT_PATH", "")
		key_path := helpers.GetEnv("IRC_DB_TLS_KEY_PATH", "")
		tls_config = helpers.CreateTLSConfig(crt_path, key_path)
	}

	db_user := helpers.GetEnv("IRC_DB_USER", "postgres")
	db_database := helpers.GetEnv("IRC_DB_DATABASE", "postgres")
	db_app_name := helpers.GetEnv("IRC_DB_APPLICATION_NAME", "zirc_server")
	db_host := fmt.Sprintf("%s:%s", helpers.GetEnv("IRC_DB_HOST", "zirc_db"), helpers.GetEnv("IRC_DB_PORT", "5432"))

	g_DB = pg.Connect(&pg.Options{
		Addr:            db_host,
		User:            db_user,
		Password:        helpers.GetEnv("IRC_DB_PASSWORD", "password"),
		Database:        db_database,
		ApplicationName: db_app_name,
		TLSConfig:       tls_config,
	})

	ctx := context.Background()
	if err := g_DB.Ping(ctx); err != nil {
		log.Info().Msgf("Could not ping database upon startup: %s", err)
	} else {
		log.Info().Msgf("Successfully connected to db. Host: %s, User: %s, Db: %s, AppName: %s", db_host, db_user, db_database, db_app_name)
	}
}

func reconnect_db_listener() {
	log.Info().Msgf("DB Reconnect Listener started")
	reconnect_wait_interval, err := strconv.Atoi(helpers.GetEnv("IRC_DB_HEALTHCHECK_INTERVAL", "3"))
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

func init() {
	enableDB, err := strconv.ParseBool(helpers.GetEnv("ZIRC_DB_ENABLED", "false"))
	if err != nil {
		log.Error().Msgf("Error parsing bool: %s", err.Error())
		enableDB = false
	}

	if enableDB {
		db_init()
		go reconnect_db_listener()
	}

	// var err error
	// plain_auth_stmt, err = g_DB.Prepare(`SELECT email, credentials from auth where username = $1::text`)
	// if err != nil {
	// 	log.Error().Msgf("Error creating plain_auth_stmt")
	// 	panic(err)
	// }

	// register_stmt, err = g_DB.Prepare(`INSERT INTO users VALUES (default, $1::text, $2::text)`)
	// if err != nil {
	// 	log.Error().Msgf("Error creating register_stmt")
	// 	panic(err)
	// }
}

type UserModel struct {
	username string
	password string
}
