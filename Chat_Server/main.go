package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/phuslu/log"
)

func main() {
	log_level := log.InfoLevel
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
		fmt.Println(err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
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
	log.Info().Msgf("Client connected: %s", remote_addr)
	recv_buf := make([]byte, 1024)

	for {
		// TODO - clear buffers so no extra data is sent?
		_, err := conn.Read(recv_buf)
		if err != nil {
			if err == io.EOF {
				log.Info().Msgf("Connection closed by client: %s", remote_addr)
			} else {
				fmt.Println(err)
			}
			return
		}

		process_message(&recv_buf)

		b := []byte("hi from irc server")
		if _, err := conn.Write(b); err != nil {
			fmt.Printf("Server got error %s\n", err)
			break
		}
	}
}

func process_message(recv_buf *[]byte) {
	trimmed_msg := bytes.Trim(bytes.TrimLeft(*recv_buf, " "), "\x00")
	log.Info().Msgf("Raw Client msg: %s", trimmed_msg)
}
