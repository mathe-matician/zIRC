package chat

import (
	"context"
	"crypto/tls"
	"fmt"
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
	enable_tls := G_Config.DB.Enable_tls

	if enable_tls {
		crt_path := G_Config.DB.Tls_cert_path
		key_path := G_Config.DB.Tls_key_path
		tls_config = helpers.CreateTLSConfig(crt_path, key_path)
	}

	db_user := G_Config.DB.User
	db_database := G_Config.DB.Database
	db_app_name := G_Config.DB.Application_name
	db_host := fmt.Sprintf("%s:%s", G_Config.DB.Host, G_Config.DB.Port)

	g_DB = pg.Connect(&pg.Options{
		Addr:            db_host,
		User:            db_user,
		Password:        G_Config.DB.Password_file,
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
	reconnect_wait_interval := G_Config.DB.Healthcheck_interval

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
	enableDB := G_Config.DB.Enabled

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
