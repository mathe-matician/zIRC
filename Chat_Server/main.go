package main

import (
	"io"
	"net"
	"strings"

	"zirc/chat"
	c "zirc/client"
	"zirc/helpers"
	rc "zirc/remote_conn"
	sm "zirc/servermanager"

	"github.com/phuslu/log"
)

const MAX_BUFFER_SIZE = 8192 // should this be a runtime configuration?

var server_dns_name = helpers.GetEnv("IRC_SERVER_DNS_NAME", "localhost")

// type RemoteConn struct {
// 	Host string
// 	Port string
// }

func main() {

	// TODO - for whatever reason overriding go's default logger
	// 		  with "github.com/phuslu/log"'s configuration doesn't show logs
	//		  within docker containers.

	// log_level := log.DebugLevel
	// level := os.Getenv("IRC_CHAT_SERVER_LOG_LEVEL")
	// if level != "" {
	// 	log_level = log.ParseLevel(level)
	// }

	// log.DefaultLogger = log.Logger{
	// 	Level:      log_level,
	// 	Caller:     1,
	// 	TimeField:  "date",
	// 	TimeFormat: "2006-01-02",
	// 	// Writer:     &log.IOWriter{os.Stderr},
	// 	Writer: &log.ConsoleWriter{
	// 		ColorOutput:    true,
	// 		QuoteString:    true,
	// 		EndWithMessage: true,
	// 	},
	// }

	// fmt.Println("IRC_DEFAULT_SERVER_NAME:", os.Getenv("IRC_DEFAULT_SERVER_NAME"))
	// fmt.Println("IRC_HOST:", os.Getenv("IRC_HOST"))
	// fmt.Println("IRC_ENABLE_TLS:", os.Getenv("IRC_ENABLE_TLS"))
	// fmt.Println("IRC_TLS_PORT:", os.Getenv("IRC_TLS_PORT"))
	// fmt.Println("IRC_VERSION:", os.Getenv("IRC_VERSION"))
	// fmt.Println("MONGODB_CHAT_SERVER_COLLECTION_NAME:", os.Getenv("MONGODB_CHAT_SERVER_COLLECTION_NAME"))

	ln, err := net.Listen("tcp", "0.0.0.0:6667")
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	server_mode := helpers.GetEnv("IRC_SERVER_ROLE", "leaf")
	err = helpers.VerifyServerMode(server_mode)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	log.Info().Msgf("Server started as %s node", server_mode)

	server_manager := sm.ServerManager{Name: server_dns_name}
	go server_manager.Run()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Error().Msgf("Error accepting connection: %s", err.Error())
			continue
		}

		go handleConnection(conn, &server_manager)
	}
}

// handleConnection is the starting point for handling a client connection
// it is intended to be run in a go routine
// and will run forever until the client connection is closed or an error is hit
func handleConnection(conn net.Conn, server_manager *sm.ServerManager) {
	defer conn.Close()

	remote_addr := conn.RemoteAddr()
	remote_ip, remote_port, err := net.SplitHostPort(remote_addr.String())
	if err != nil {
		log.Error().Str("remote_addr", remote_addr.String()).Msgf("Error splitting remote addr: %s", err.Error())
	}
	// TODO - resolve DNS name here for additional checks / verification
	// e.g. w/ servers and compare to server list
	remote_conn := rc.NewRemoteConn("", remote_ip, remote_port)
	client, session_timestamp, err := c.NewClient("", "", remote_conn)
	if err != nil {
		log.Error().EmbedObject(client).Msgf(err.Error())
		return
	}

	// add the client to the global client list
	server_manager.ClientList = append(server_manager.ClientList, client)

	log.Info().EmbedObject(client).Msgf("Client connected at %s", *session_timestamp)

	recv_buf := make([]byte, MAX_BUFFER_SIZE)

	for {
		// TODO - clear buffers so no extra data is sent?
		// TODO - send periodic PING commands
		//		  if no PONG is received, terminate the connection
		//		  used to determine dead connections
		_, err := conn.Read(recv_buf)
		if err != nil {
			if err == io.EOF {
				end_timestamp, err := client.SetSessionEndTimestamp()
				if err != nil {
					log.Error().EmbedObject(client).Msg(err.Error())
				}
				log.Info().EmbedObject(client).Msgf("Client disconnected: %s", *end_timestamp)
				// TODO - write session duration as a metric / possible analysis
			} else {
				log.Error().EmbedObject(client).Msgf("Error reading data from connection: %s", err.Error())
			}
			return
		}

		response := chat.ProcessMessage(&recv_buf, client, server_manager)

		if len(response) == 0 {
			// e.g. sometimes the server doesn't send anything back to the client
			// 		as in the case of correct password via PASS
			continue
		}

		if _, err := conn.Write(response); err != nil {
			log.Error().EmbedObject(client).Msgf("Error writing to client: %s", err.Error())
			break
		}

		if strings.Contains(string(response), "ERROR") {
			log.Error().EmbedObject(client).Msgf("Critical error occurred. Closing client connection: %s", response)
			conn.Close()
			break
		}
	}
}
