package main

import (
	"io"
	"net"

	"zirc/chat"
	"zirc/helpers"
	rc "zirc/remote_conn"
	sm "zirc/servermanager"

	"github.com/google/uuid"
	"github.com/phuslu/log"
)

const MAX_BUFFER_SIZE = 4096

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

	server_manager := sm.ServerManager{}
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
	remote_conn := rc.RemoteConn{Host: "", Ip: remote_ip, Port: remote_port}
	uuid, err := uuid.NewV7()

	if err != nil {
		log.Error().EmbedObject(&remote_conn).Msgf("Error generating uuid %s", err.Error())
		return
	}

	timestamp_s, timestamp_ns := uuid.Time().UnixTime()
	log.Info().EmbedObject(&remote_conn).Msgf("Client connected at %d.%d", timestamp_s, timestamp_ns)

	recv_buf := make([]byte, MAX_BUFFER_SIZE)

	for {
		// TODO - clear buffers so no extra data is sent?
		// TODO - send periodic PING commands
		//		  if no PONG is received, terminate the connection
		//		  used to determine dead connections
		_, err := conn.Read(recv_buf)
		if err != nil {
			if err == io.EOF {
				log.Info().EmbedObject(&remote_conn).Msg("Client disconnected")
			} else {
				log.Error().EmbedObject(&remote_conn).Msgf("Error reading data from connection: %s", err.Error())
			}
			return
		}

		response := chat.ProcessMessage(&recv_buf, server_manager)
		// response := []byte("hi from IRC...!")

		if _, err := conn.Write(response); err != nil {
			log.Error().EmbedObject(&remote_conn).Msgf("Error writing to client: %s", err.Error())
			break
		}
	}
}
