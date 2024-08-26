package servermanager

import (
	"fmt"

	c "zirc/client"

	"github.com/phuslu/log"
)

type ServerManager struct {
	Name       string
	recv       chan string
	ClientList []c.Client // only contains registered clients
	ServerList []string
	// send        chan string
}

// IsClientRegistered checks to see if a particular client's session_id
// is in the ServerManager's ClientList. ClientList only contains registered clients,
// so if a session id exists at all, that means the client is registered
func (sm *ServerManager) IsClientRegistered(session_id string) bool {
	for _, c := range sm.ClientList {
		if c.SessionId() == session_id {
			return true
		}
	}
	return false
}

func (sm *ServerManager) Run() {
	log.Info().Msg("ServerManager started")
	// setup global recv and send channels
	// these are used for each connection to relay information back to the server to then send
	// to other clients or servers
	sm.recv = make(chan string)
	// sm.send = make(chan string) // what was the send chan for again?? I don't think it is necessary

	for {
		// TODO
		// if commands are of certain types, the ServerManager will need to send messages to
		// send to all clients in a channel
		// for example:
		// 		- JOIN (user1 has joined the channel)
		//		- NOTIFY
		//		... etc
		msg1 := <-sm.recv
		fmt.Println(msg1)

		// select {
		// case msg1 := <-sm.recv:

		// 	fmt.Println(msg1)
		// 	// if one of these commands...
		// 	// for i, client := range sm.client_list {
		// 	// 	// o(n) - not the best, but probably fine for now

		// 	// }

		// 	// case msg2 := <-sm.send:
		// 	// 	fmt.Println(msg2)
		// }
	}
}
