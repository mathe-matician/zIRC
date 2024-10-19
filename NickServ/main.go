package nickserv

import "github.com/phuslu/log"

func connect() {

}

func main() {
	log.Info().Msg("Starting NickServ")
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
