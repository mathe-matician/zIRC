package servermanager

import (
	"fmt"

	c "zirc/client"

	"github.com/phuslu/log"
)

type ServerManager struct {
	name        string
	recv        chan string
	client_list []c.Client
	server_list []string
	// send        chan string
}

func (sm *ServerManager) Run() {
	log.Info().Msg("ServerManager started")
	// setup global recv and send channels
	// these are used for each connection to relay information back to the server to then send
	// to other clients or servers
	sm.recv = make(chan string)
	// sm.send = make(chan string) // what was the send chan for again?? I don't think it is necessary

	for {
		select {
		case msg1 := <-sm.recv:
			// TODO
			// if commands are of certain types, the ServerManager will need to send messages to
			// send to all clients in a channel
			// for example:
			// 		- JOIN (user1 has joined the channel)
			//		- NOTIFY
			//		... etc
			fmt.Println(msg1)
			// if one of these commands...
			// for i, client := range sm.client_list {
			// 	// o(n) - not the best, but probably fine for now

			// }

			// case msg2 := <-sm.send:
			// 	fmt.Println(msg2)
		}
	}
}
