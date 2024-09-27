package chat

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"zirc/helpers"

	"github.com/phuslu/log"
)

var server_manager_commands = map[string]string{
	"PASS":   "",
	"SERVER": "",
	"CAP":    "",
	"SJOIN":  "",
	"PING":   "",
	"PONG":   "",
	"SQUIT":  "",
}

type ServerManager struct {
	Addr              string
	serverIpWhitelist map[string]string
}

func NewServerManager() *ServerManager {
	return &ServerManager{
		Addr: helpers.GetEnv("IRC_SERVER_MANAGER_PORT", "7000"),
	}
}

func (sm *ServerManager) UpdateServerIpWhiteList(ip, args string) {
	sm.serverIpWhitelist[ip] = args
}

// Allows this server to Join an existing IRC network
func (sm *ServerManager) Join() {

}

// Checks to see whether the server trying to connect to this network
// has entered in a valid password via the PASS command
func (sm *ServerManager) ValidS2SPassword(password string) bool {
	// check db first
	// if can't connect to db check locally for conf file or env variable
	// if all fails, then don't let join

	// last fallback is env var
	if helpers.GetEnv("IRC_S2S_PASSWORD", "") == password {
		return true
	}
	return false
}

func (sm *ServerManager) Run() {
	// listens on port 7000 for servers joining the irc network
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()
	ln := GetTCPListener(
		helpers.GetEnv("IRC_S2S_ENABLE_TLS", "false"),
		helpers.GetEnv("IRC_S2S_TLS_CERT_PATH", "./.tls/s2s.crt"),
		helpers.GetEnv("IRC_S2S_TLS_KEY_PATH", "./.tls/s2s.key"),
		helpers.GetEnv("IRC_S2S_TLS_PORT", "7001"),
		helpers.GetEnv("IRC_S2S_PORT", "7000"),
	)

	log.Info().Msgf("ServerManager listening: %s", sm.Addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Error().Msgf("Error accepting connection: %s", err.Error())
			continue
		}

		if conn.RemoteAddr().Network() != "tcp" {
			conn.Close()
			continue
		}

		if !sm.whiteListedServerIp(conn.RemoteAddr().String()) {
			conn.Close()
			continue
		}

		go sm.handleConnection(&conn)
	}
}

// Checks to see if the connected IP is an expected server IP
// i.e. do we expect this IP to be communicating with the ServerManager
func (sm *ServerManager) whiteListedServerIp(ip string) bool {
	server_ip, _, err := net.SplitHostPort(ip)
	if err != nil {
		log.Error().Msgf("ServerManager Error extracting IP: %s", err)
		return false
	}

	if _, ok := sm.serverIpWhitelist[server_ip]; !ok {
		log.Error().Msgf("ServerManager not a whitelisted server ip!: %s", server_ip)
		return false
	}
	return true
}

func (sm *ServerManager) handleConnection(conn *net.Conn) {
	defer (*conn).Close()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ServerManager recovered from panic:", r)
		}
	}()
	remote_addr := (*conn).RemoteAddr()
	log.Info().Msgf("Handling ServerManager connection: %s", remote_addr)

	var max_buffer_size int
	var err error
	// TODO
	// does this size differ for S2S?
	max_buffer_size, err = strconv.Atoi(helpers.GetEnv("IRC_S2S_MAX_BUFFER_SIZE", "262144"))
	if err != nil {
		log.Error().Msgf(err.Error())
		max_buffer_size = 262144
	}

	for {
		if !sm.whiteListedServerIp((*conn).RemoteAddr().String()) {
			// need to always check if this IP is still valid as we handle this connection
			// as we could update the whitelist mid connection
			break
		}

		// block on read until the buffer has at least 1 byte.
		// just a hacky way for this to block as Read() doesn't block on its own
		recv_buf := make([]byte, max_buffer_size)
		_, err := io.ReadAtLeast((*conn), recv_buf, 1)
		if err != nil {
			if err == io.EOF {
				log.Info().Msgf("Server %s disconnected", remote_addr)
			}
			return
		}

		if _, err := (*conn).Write([]byte("")); err != nil {
			break
		}
	}
}
