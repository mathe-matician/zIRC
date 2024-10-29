package nickserv

import (
	"io"
	"net"
	"regexp"
	"strconv"
	"time"

	"github.com/phuslu/log"
)

// captures until first space and has two capture groups, what comes before the space and what comes after
// e.g. `#general hello world` would contain two capture groups: ((#general), (hello world))
// meant to be used with FindStringSubmatch(str)
var cmd_re = regexp.MustCompile(`^(\S+)(.*)`)

type NickServ struct {
	ConnListener *net.Listener
}

func NewNickServ() *NickServ {
	enable_tls := GetEnv("NICKSERV_ENABLE_TLS", "false")
	port := GetEnv("NICKSERV_PORT", "29999")

	listener := GetTCPListener(
		enable_tls,
		GetEnv("NICKSERV_TLS_CERT_PATH", "./.tls/nickserv.crt"),
		GetEnv("NICKSERV_TLS_KEY_PATH", "./.tls/nickserv.key"),
		GetEnv("NICKSERV_TLS_PORT", "30000"),
		port,
	)

	log.Info().Msgf("Creating new NickServ instance: port: %s, tls_enabled: %s", port, enable_tls)

	return &NickServ{
		ConnListener: &listener,
	}
}

func (ns *NickServ) HandleConnection(conn *net.Conn) {
	defer (*conn).Close()
	remote_addr := (*conn).RemoteAddr()

	// remote_ip, remote_port, err := net.SplitHostPort(remote_addr.String())
	// if err != nil {
	// 	log.Error().Str("remote_addr", remote_addr.String()).Msgf("Error splitting remote addr: %s", err.Error())
	// }

	log.Info().Msgf("Client connected: %s", remote_addr.String())

	max_buffer_size, err := strconv.Atoi(GetEnv("NICKSERV_MAX_BUFFER_SIZE", "8192"))
	if err != nil {
		log.Error().Msgf(err.Error())
		max_buffer_size = 8192
	}

	for {
		// block on read until the buffer has at least 1 byte.
		// just a hacky way for this to block as Read() doesn't block on its own
		recv_buf := make([]byte, max_buffer_size)
		_, err := io.ReadAtLeast((*conn), recv_buf, 1)
		if err != nil {
			if err == io.EOF {
				// end_timestamp, err := client.SetSessionEndTimestamp()
				// if err != nil {
				// 	log.Error().EmbedObject(client).Msg(err.Error())
				// }
				log.Info().Msgf("Client %s disconnected", remote_addr.String())
				// TODO - remove client state from MessageManager!!!
				// TODO - write session duration as a metric / possible analysis

				// TODO - remove client from all channels they are on
			} else {
				log.Error().Msgf("Error reading data from connection: %s", err.Error())
			}

			//if err == io.ErrShortBuffer
			return
		}

		msg := string(recv_buf)

		log.Debug().Msgf("Client sent: %s", msg)

		if len(msg) <= 0 {
			return
		}

		// TODO
		// timeout should be based on specific commands
		// e.g. if I'm an admin, I probably don't want my connection timing out
		// if there is back and forth type commands
		// that is, ONLY IF those types of commands exist
		conn_timeout_duration, err := strconv.Atoi(GetEnv("NICKSERV_CONN_TIMEOUT_DURATION", "2"))
		if err != nil {
			log.Error().Msgf("Error parsing conn timeout duration env var")
			panic(err)
		}

		// TODO
		// (determine IF this needs to be conditionally applied)
		// Since interactions with NickServ shouldn't be long running
		// set a dealine for read interactions on this socket
		// This read deadline is the timeout between socket reads
		// e.g. a client is performing some command
		//		IF there should be interactions back and forth between client and server
		//		have each client interaction timeout at max after X seconds
		//		if the client sends a msg within the window, that deadline renews
		(*conn).SetReadDeadline(time.Now().Add(time.Duration(conn_timeout_duration) * time.Minute))

		// Get params sent in msg
		split_msg := cmd_re.FindStringSubmatch(msg)
		log.Debug().Msgf("split_msg: %s", split_msg)
		cmd := split_msg[0]

		// get func from command map
		fn, ok := command_map[cmd]
		if !ok {
			log.Error().Msgf("Not a valid NickServ command: %s", cmd)
			// TODO
			// should this return something to the client?
			break
		}

		params := split_msg[1]
		response := fn(params)

		if _, err := (*conn).Write([]byte(response)); err != nil {
			log.Error().Msgf("Error writing to client: %s", err.Error())
			break
		}
	}
}
