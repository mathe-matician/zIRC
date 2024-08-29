package tests

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type MockClient struct{}

func (mc *MockClient) Send(addr string, data string, timeout bool) string {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer conn.Close()

	d := []byte(data)
	_, err = conn.Write(d)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	reply := make([]byte, 1024)

	fmt.Println("Waiting for server reply...")

	if timeout {
		ok := make(chan bool)

		go readServerResponse(&conn, &reply, ok)

		select {
		case <-ok:
			fmt.Println("Client received response from server")
		case <-time.After(2 * time.Second):
			fmt.Println("Expected client timeout!")
			return ""
		}
	} else {
		readServerResponse(&conn, &reply, nil)
	}

	// remove extra null characters from 1024 buffer
	reply = bytes.Trim(reply, "\x00")

	return string(reply)
}

func readServerResponse(conn *net.Conn, buffer *[]byte, ok chan bool) {
	_, err := (*conn).Read(*buffer)
	if err != nil {
		if !strings.Contains(err.Error(), "use of closed network connection") {
			fmt.Println("Write to server failed: ", err.Error())
			os.Exit(1)
		}
	}
	if ok != nil {
		ok <- true
	}
}
