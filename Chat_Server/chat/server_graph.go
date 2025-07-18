package chat

type ServerGraph struct {
	Graph []*IrcServer
}

func NewServerGraph() *ServerGraph {
	graph := make([]*IrcServer, 0)

	return &ServerGraph{
		Graph: graph,
	}
}

// Insert inserts an IrcServer into the server graph DFS based on some parent
// If parent == nil then just append the server.
// This generally happens when the first server is added
func (sg *ServerGraph) Insert(server *IrcServer, parent *IrcServer) {
	if parent == nil {
		sg.Graph = append(sg.Graph, server)
		return
	}

	stack := make([]*IrcServer, 0)
	stack = append(stack, server)

	for len(stack) != 0 {
		stackLen := len(stack)
		last := stackLen - 1

		currentServer := stack[last]
		stack = stack[:last]

		if currentServer == parent {
			parent._ServerGraph.Graph = append(parent._ServerGraph.Graph, server)
			return
		}

		for _, svr := range currentServer._ServerGraph.Graph {
			stack = append(stack, svr)
		}
	}
}

func (sg *ServerGraph) Remove(server *IrcServer) {

}
