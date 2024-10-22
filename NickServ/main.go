package main

import (
	src "nickserv/src"

	"github.com/phuslu/log"
)

var NS *src.NickServ

func main() {
	log.Info().Msg("Starting NickServ")
	src.DB_pkg_init()
	NS = src.NewNickServ()

	for {
		conn, err := (*NS.ConnListener).Accept()
		if err != nil {
			log.Error().Msgf("Error accepting connection: %s", err.Error())
			continue
		}

		if conn.RemoteAddr().Network() != "tcp" {
			conn.Close()
			continue
		}

		go NS.HandleConnection(&conn)
	}
}
