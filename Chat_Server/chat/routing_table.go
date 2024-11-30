package chat

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

// handles msgs
func (rt *RoutingTable) EventListener() {

}

func (rt *RoutingTable) Update() {

}
