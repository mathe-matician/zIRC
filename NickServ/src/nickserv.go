package nickserv

import (
	"io"
	"net"
	"strconv"

	"github.com/phuslu/log"
)

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
			if conn != nil {
				(*conn).Close()
			}
			return
		}

		msg := string(recv_buf)

		log.Debug().Msgf("Client sent: %s", msg)

		if len(msg) <= 0 {
			if conn != nil {
				(*conn).Close()
			}
			return
		}
	}
}
