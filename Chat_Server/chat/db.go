package chat

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"time"
	"zirc/helpers"

	"github.com/go-pg/pg"
	"github.com/phuslu/log"
)

var g_DB *pg.DB

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

	g_DB = pg.Connect(&pg.Options{
		Addr:            fmt.Sprintf("%s:%s", helpers.GetEnv("IRC_DB_HOST", "zirc_db"), helpers.GetEnv("IRC_DB_PORT", "5432")),
		User:            helpers.GetEnv("IRC_DB_USER", "postgres"),
		Password:        helpers.GetEnv("IRC_DB_PASSWORD", "password"),
		Database:        helpers.GetEnv("IRC_DB_DATABASE", "postgres"),
		ApplicationName: helpers.GetEnv("IRC_DB_APPLICATION_NAME", "zirc_server"),
		TLSConfig:       tls_config,
	})
}

func init() {
	db_init()
	go reconnect_db_listener()
}

func reconnect_db_listener() {
	reconnect_wait_interval, err := strconv.Atoi(helpers.GetEnv("IRC_DB_RECONNECT_WAIT_INTERVAL", "3"))
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
