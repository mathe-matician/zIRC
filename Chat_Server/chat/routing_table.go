package chat

import (
	"errors"

	"github.com/phuslu/log"
)

type RoutingAction struct {
}

// e.g.
//    D - B - C
//        |
//    this_server
//
// Destination	Next Hop
// B			Direct
// C			B
// D			B

// Servers map contains the name of the server which maps to its next hop
type RoutingTable struct {
	Servers map[string]string
	Recv    chan RoutingAction
}

func NewRoutingTable() *RoutingTable {
	return &RoutingTable{
		make(map[string]string),
		make(chan RoutingAction),
	}
}

// GetServer iterates through server map
//
// continue to search routing table until you find the root server which connects
// to this server
// worst case o(n)
// best o(1)
// TODO
// get other features to determine whether to send along a specific path
// e.g. hop count
// e.g. latency etc
func (rt *RoutingTable) GetServer(server string) (string, int, error) {
	next := server
	var ok bool
	hopcount := 0
	// TODO
	// this doesn't account for situations where
	// you can't reach the server, but it is still part of the network
	// or does it?
	log.Debug().Msgf("Routing Table: getting server: %s", server)
	for next != "direct" {
		next, ok = rt.Servers[next]
		if !ok {
			return "", -1, errors.New("server not in routing table")
		}
		hopcount++
	}
	return next, hopcount, nil
}

// handles msgs
func (rt *RoutingTable) EventListener() {

}

func (rt *RoutingTable) Update() {

}
