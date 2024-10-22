package main

import (
	nickserv "nickserv/src"

	"github.com/phuslu/log"
)

func main() {
	log.Info().Msg("Starting NickServ")
	nickserv.DB_pkg_init()
	for {
		// conn, err := (*is.Listener).Accept()
		// if err != nil {
		// 	log.Error().Msgf("Error accepting connection: %s", err.Error())
		// 	continue
		// }

		// if conn.RemoteAddr().Network() != "tcp" {
		// 	conn.Close()
		// 	continue
		// }

		// go is.handleConnection(&conn)
	}
}
