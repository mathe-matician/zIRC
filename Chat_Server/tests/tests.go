package tests

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"github.com/phuslu/log"
)

func init() {
	// used for random NICK and USER name generation
	rand.Seed(time.Now().UnixNano())
}

type MockMsg struct {
	IrcMsg              string
	WaitForResponseCode string
	Timeout             bool
}

type MockResponse struct {
	Response     string
	CheckSuccess bool
	Success      bool
}

type MockClient struct {
	Connection net.Conn
	Addr       string
	SendChan   chan MockMsg      // contains message to send to irc server
	Recv       chan MockResponse // the chan our mock client can send server responses which we can then check in our test
}

func NewMockClient(addr string, recv chan MockResponse) *MockClient {
	if addr == "" {
		panic("missing addr arg for new mock client")
	}
	sendChan := make(chan MockMsg, 30)

	return &MockClient{
		Connection: nil,
		Addr:       addr,
		SendChan:   sendChan,
		Recv:       recv,
	}
}

func (mc *MockClient) Send(msg, waitForResponseCode string, timeout bool) {
	// time.Sleep(3 * time.Second)
	mc.SendChan <- MockMsg{msg, waitForResponseCode, timeout}
}

func WaitUntilResponseCode(recv chan string, response_code string) {
	fmt.Printf("Waiting until recv response code: %s\n", response_code)
	for {
		select {
		case res := <-recv:
			log.Debug().Msgf("WaitUntilResponseCode recvd: %s", res)
			// fmt.Printf("Recvd %s\n", res)
			res_split := strings.Split(res, "\r\n")
			log.Debug().Msgf("Res split: %v", res_split)
			for _, v := range res_split {
				if strings.Contains(v, response_code) {
					log.Debug().Msgf("Response contains response code: %s", response_code)
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
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer mc.Connection.Close()

	bufReader := bufio.NewReader(mc.Connection)
	for {
		select {
		case msg := <-mc.SendChan:
			_, err = mc.Connection.Write([]byte(msg.IrcMsg))
			if err != nil {
				fmt.Println("Error sending message:", err)
				os.Exit(1)
			}
			time.Sleep(3 * time.Second)

			if msg.Timeout {
				ok := make(chan bool)
				go func() {
					response := readServerResponse(mc.Connection, ok)
					fmt.Println("Response:", response)
					mc.Recv <- MockResponse{response, false, false}
				}()
				select {
				case <-ok:
				case <-time.After(2 * time.Second):
					fmt.Println("Expected timeout!")
				}
			} else {
				response, _ := bufReader.ReadString('\n')
				fmt.Println("Client Received:", response)
				mc.Recv <- MockResponse{response, false, false}
			}

		case <-time.After(10 * time.Second):
			fmt.Println("No messages to process. Exiting.")
			return
		}
	}
}

func (mc *MockClient) Runold() {
	var err error
	mc.Connection, err = net.Dial("tcp", mc.Addr)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer mc.Connection.Close()
	mc.Recv <- MockResponse{}

	for {
		msg := <-mc.SendChan
		d := []byte(msg.IrcMsg)
		reply := make([]byte, 1024)
		_, err = mc.Connection.Write(d)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		fmt.Println("Waiting for server reply...")

		response := ""

		if msg.Timeout {
			ok := make(chan bool)

			go readServerResponse(mc.Connection, ok)

			select {
			case <-ok:
				fmt.Println("Client received response from server")
			case <-time.After(2 * time.Second):
				fmt.Println("Expected client timeout!")
				return
			}
		} else {
			response = readServerResponse(mc.Connection, nil)
			log.Debug().Msgf("RESPONSE:::: %v", response)
		}

		// remove extra null characters from 1024 buffer
		// log.Debug().Msgf("REPLY:::: %v", string(reply))
		reply = bytes.Trim(reply, "\x00")
		if msg.WaitForResponseCode != "" {
			log.Debug().Msgf("waiting... recvd: %s", reply)
			res_split := strings.Split(string(reply), "\r\n")
			log.Debug().Msgf("Res split: %v", res_split)
			for _, v := range res_split {
				if strings.Contains(v, msg.WaitForResponseCode) {
					log.Debug().Msgf("Response contains response code: %s", msg.WaitForResponseCode)
					mc.Recv <- MockResponse{"", true, true}
				}
			}
		} else {
			// mc.Recv <- MockResponse{string(reply), false, false}
			mc.Recv <- MockResponse{response, false, false}
		}
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

		go readServerResponse(mc.Connection, ok)

		select {
		case <-ok:
			fmt.Println("Client received response from server")
		case <-time.After(2 * time.Second):
			fmt.Println("Expected client timeout!")
			return ""
		}
	} else {
		readServerResponse(mc.Connection, nil)
	}

	// remove extra null characters from 1024 buffer
	reply = bytes.Trim(reply, "\x00")

	return string(reply)
}

func readServerResponse(conn net.Conn, ok chan bool) string {
	reply := make([]byte, 1024)
	_, err := conn.Read(reply)
	if err != nil {
		if !strings.Contains(err.Error(), "use of closed network connection") {
			fmt.Println("Write to server failed: ", err.Error())
			os.Exit(1)
		}
	}

	if ok != nil {
		ok <- true
	}

	return string(bytes.Trim(reply, "\x00"))
}
