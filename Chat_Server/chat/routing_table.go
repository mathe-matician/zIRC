package chat

type RoutingAction struct {
}

type RoutingTable struct {
	Neighbors []string // should this be map?
	Recv      chan RoutingAction
}

func NewRoutingTable() *RoutingTable {
	return &RoutingTable{}
}

// handles msgs
func (rt *RoutingTable) EventListener() {

}

func (rt *RoutingTable) Update() {

}
