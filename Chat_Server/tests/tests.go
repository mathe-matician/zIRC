package tests

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

func init() {
	// used for random NICK and USER name generation
	rand.Seed(time.Now().UnixNano())
}

type MockMsg struct {
	IrcMsg  string
	Timeout bool
}

type MockClient struct {
	Connection net.Conn
	Addr       string
	SendChan   chan MockMsg // contains message to send to irc server
	Recv       chan string  // the chan our mock client can send server responses which we can then check in our test
}

func NewMockClient(addr string, recv chan string) *MockClient {
	if addr == "" {
		panic("missing addr arg for new mock client")
	}
	sendChan := make(chan MockMsg)

	return &MockClient{
		Connection: nil,
		Addr:       addr,
		SendChan:   sendChan,
		Recv:       recv,
	}
}

func (mc *MockClient) Send(msg string, timeout bool) {
	mc.SendChan <- MockMsg{msg, timeout}
}

func WaitUntilResponseCode(recv chan string, response_code string) {
	fmt.Printf("Waiting until recv response code: %s\n", response_code)
	for {
		select {
		case res := <-recv:
			fmt.Printf("Recvd %s\n", res)
			res_split := strings.Split(res, "\r\n")
			for _, v := range res_split {
				if strings.Contains(v, response_code) {
					return
				}
			}
			// case <-time.After(10 * time.Second):
			// 	panic("waituntil timed out unexpectidly")
		}
	}
}

func (mc *MockClient) Run() {
	var err error
	mc.Connection, err = net.Dial("tcp", mc.Addr)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer mc.Connection.Close()

	for {
		msg := <-mc.SendChan
		d := []byte(msg.IrcMsg)
		_, err = mc.Connection.Write(d)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		reply := make([]byte, 1024)
		fmt.Println("Waiting for server reply...")

		if msg.Timeout {
			ok := make(chan bool)

			go readServerResponse(mc.Connection, &reply, ok)

			select {
			case <-ok:
				fmt.Println("Client received response from server")
			case <-time.After(2 * time.Second):
				fmt.Println("Expected client timeout!")
				return
			}
		} else {
			readServerResponse(mc.Connection, &reply, nil)
		}

		// remove extra null characters from 1024 buffer
		reply = bytes.Trim(reply, "\x00")
		mc.Recv <- string(reply)
	}
}

// https://stackoverflow.com/a/22892986
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RegisterMsg(addr, nick, user string) string {
	if nick == "" {
		nick = randSeq(5)
	}

	if user == "" {
		user = randSeq(5)
	}
	return fmt.Sprintf("NICK %s\r\nUSER %s\r\n", nick, user)
}

func (mc *MockClient) OLDClientSendChan(addr string, data string, timeout bool) string {
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

		go readServerResponse(mc.Connection, &reply, ok)

		select {
		case <-ok:
			fmt.Println("Client received response from server")
		case <-time.After(2 * time.Second):
			fmt.Println("Expected client timeout!")
			return ""
		}
	} else {
		readServerResponse(mc.Connection, &reply, nil)
	}

	// remove extra null characters from 1024 buffer
	reply = bytes.Trim(reply, "\x00")

	return string(reply)
}

func readServerResponse(conn net.Conn, buffer *[]byte, ok chan bool) {
	_, err := conn.Read(*buffer)
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
