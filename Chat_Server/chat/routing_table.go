package chat

import "errors"

type RoutingAction struct {
}

// e.g.
// Destination	Next Hop
// B			Direct
// C			B
// D			B

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

// continue to search routing table until you find the root server which connects
// to this server
// worst case o(n)
// best o(1)
// TODO
// get other features to determine whether to send along a specific path
// e.g. hop count
// e.g. latency etc
func (rt *RoutingTable) GetServer(server string) (string, error) {
	next := server
	var ok bool
	for next != "direct" {
		next, ok = rt.Servers[next]
		if !ok {
			return "", errors.New("server not in routing table")
		}
	}
	return next, nil
}

// handles msgs
func (rt *RoutingTable) EventListener() {

}

func (rt *RoutingTable) Update() {

}
