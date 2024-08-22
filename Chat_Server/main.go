package main

import (
	"fmt"
	"net"
)

func main() {
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
	fmt.Printf("Client connected: %s\n", conn.RemoteAddr())
	server_buf := make([]byte, 1024)
	// client_buf := make([]byte, 1024)

	for {
		// TODO - clear buffers so no extra data is sent?
		_, err := conn.Read(server_buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Client sent: %s\n", server_buf)
		b := []byte("hi from irc server")
		if _, err := conn.Write(b); err != nil {
			fmt.Printf("Server got error %s\n", err)
			break
		}
	}
}
