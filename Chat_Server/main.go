package main

import (
	"io"
	"net"
	"os"

	"zirc/chat"

	"github.com/phuslu/log"
)

type RemoteConn struct {
	Host string
	Port string
}

func (c *RemoteConn) MarshalObject(e *log.Entry) {
	e.Str("host", c.Host).Str("port", c.Port)
}

func main() {
	log_level := log.DebugLevel
	level := os.Getenv("IRC_CHAT_SERVER_LOG_LEVEL")
	if level != "" {
		log_level = log.ParseLevel(level)
	}

	log.DefaultLogger = log.Logger{
		Level:      log_level,
		Caller:     1,
		TimeField:  "date",
		TimeFormat: "2006-01-02",
		Writer:     &log.IOWriter{os.Stdout},
	}

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

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Error().Msgf("Error accepting connection: %s", err.Error())
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// TODO - send periodic PING commands
	//		  if no PONG is received, terminate the connection
	//		  used to determine dead connections
	defer conn.Close()

	remote_addr := conn.RemoteAddr()
	remote_host, remote_port, err := net.SplitHostPort(remote_addr.String())
	if err != nil {
		log.Error().Str("remote_addr", remote_addr.String()).Msgf("Error splitting remote addr: %s", err.Error())
	}
	remote_conn := RemoteConn{Host: remote_host, Port: remote_port}
	log.Info().EmbedObject(&remote_conn).Msg("Client connected")

	recv_buf := make([]byte, 1024)

	for {
		// TODO - clear buffers so no extra data is sent?
		_, err := conn.Read(recv_buf)
		if err != nil {
			if err == io.EOF {
				log.Info().EmbedObject(&remote_conn).Msg("Client disconnected")
			} else {
				log.Error().EmbedObject(&remote_conn).Msgf("Error reading data from connection: %s", err.Error())
			}
			return
		}

		response := chat.ProcessMessage(&recv_buf)

		if _, err := conn.Write(response); err != nil {
			log.Error().EmbedObject(&remote_conn).Msgf("Error writing to client: %s", err.Error())
			break
		}
	}
}
