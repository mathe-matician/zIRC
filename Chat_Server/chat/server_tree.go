package chat

import (
	"fmt"
	"net"
	"zirc/helpers"

	"github.com/phuslu/log"
)

/*
zirc_chat_server internal tree example:

Server topology:
zirc_chat_server <- zirc_zoo <- zirc_fake1 <- zirc_fake3

	                	^
		                |
					zirc_fake2

	{
		"zirc_fake3": {
			Parent: "zirc_fake1",
			DirectlyConnected: false,
			Server: IrcServer(zirc_fake3)
		},
		"zirc_fake1": {
			Parent: "zirc_zoo",
			DirectlyConnected: false,
			Server: IrcServer(zirc_fake1)
		},
		"zirc_zoo":   {
			Parent: "zirc_chat_server",
			DirectlyConnected: true,
			Server: IrcServer(zirc_zoo)
		},
		"zirc_fake2": {
			Parent: "zirc_zoo",
			DirectlyConnected: false,
			Server: IrcServer(zirc_fake2)
		},
	}

	Pending map[net.Conn]*ServerNode:
		During daemon startup, you can specify a list of already known servers to connect to.
		During this initial handshake, the connecting server won't know the receiving server's topology, name, or SID.
		So there needs to be a temporary pending server map, as the only identifiable piece of information
		we have about it is the socket (net.Conn).
		Once we complete the handshake, this should be removed from the Pending map.
*/
type ServerTree struct {
	Tree    map[string]*ServerNode
	Pending map[net.Conn]*ServerNode
}

func NewServerTree(tree map[string]*ServerNode, pending map[net.Conn]*ServerNode) *ServerTree {
	if tree == nil {
		tree = make(map[string]*ServerNode)
	}

	if pending == nil {
		pending = make(map[net.Conn]*ServerNode)
	}

	return &ServerTree{
		Tree:    tree,
		Pending: pending,
	}
}

// Remove cleans removes the specified server from the tree
// as well as all of its downstream servers
func (st *ServerTree) Remove(server *ServerNode) {
	defer server.Conn.Conn.Close()

	for _, s := range st.Tree[server.Name].Servers {
		if s.Conn != nil {
			s.Conn.Conn.Close()
		}
		// remove all of this server's downstream servers from the main tree map
		delete(st.Tree, s.Name)
	}
	// finally remove the server from the main tree map
	delete(st.Tree, server.Name)
}

func (st *ServerTree) GetServerByConn(c net.Conn) *ServerNode {
	for _, s := range st.Tree {
		if s.IsMe {
			// our conn will always be nil
			// as we would never have a self conn
			continue
		}

		if s == nil {
			st.Remove(s)
			continue
		}

		if s.Conn.Conn == c {
			return s
		}
	}

	return nil
}

// ServerBurstSend traverses the server tree from a given root node
// Uses BFS as we usually traverse to send messages in topological order
// i.e. we should always be sending to the upstream parent servers first then to the downstream child servers
//
// root *ServerNode: Where we want to start the traversal from
// c net.Conn: Where to send the SERVER commands
func (st *ServerTree) ServerBurstSend(root *ServerNode, c net.Conn) {
	if st.Tree[root.Name] == nil {
		return
	}

	log.Debug().Msg("BURSTing my existing servers")

	q := make([]*ServerNode, 0)
	q = append(q, st.Tree[root.Name])

	for len(q) != 0 {
		log.Debug().Msgf("q: %v", q)
		currentNode := q[0]
		q = q[1:]

		if currentNode == nil {
			log.Error().Msg("currentNode is nil should probably be cleaned up!")
			continue
		}
		// always add +1 to hopcount as from the receiving server's perspective the hop will have in increased by 1
		serverCmd := fmt.Sprintf("SERVER %s %d :%s %s", currentNode.Name, currentNode.HopCount, currentNode.Description, "\r\n")

		log.Debug().Msgf("Bursting: %s", serverCmd)

		_, err := c.Write([]byte(serverCmd))
		if helpers.IsNetConnClosedErr(err) {
			log.Error().Msgf("ServerTree.Traverse: trying to write to closed connection for %s: %s", currentNode.Name, err.Error())
			st.Remove(currentNode)
			continue
		}

		for _, svr := range currentNode.Servers {
			if svr.Conn.Conn == c {
				// don't include the server we are bursting to
				// to avoid loop
				continue
			}
			q = append(q, svr)
		}
	}
}

// Send traverses the tree propegating the message to all servers
// func (st *ServerTree) Send(c net.Conn, tasks []*Task) []*ServerNode {
// 	for _, task := range tasks {
// 		_, err := c.Write([]byte(task.Task))
// 		if helpers.IsNetConnClosedErr(err) {
// 			// TODO
// 			// clean up closed connection!!
// 			log.Error().Msgf("SERVER(s2s_comm): trying to write to closed Client Conn %s", err.Error())
// 			return nil
// 		}
// 	}

// 	return nil
// }

// Insert adds server as an entry into the ServerTree map
// func (st *ServerTree) Insert(sid, name string, directlyConnected, isMe bool, hopCount int, conn *ServerConn, parent *ServerNode, server *IrcServer) {
// 	// TODO
// 	// what happens if this server already exists?

// 	// TODO
// 	// maybe check Conn.State?

// }

// Get returns the server tree in order where each parent appears before its child nodes
// For example, given a server tree of A - B - C - D where the server running this function is B after A connects to it
// It would return {B, C, D}
//
// Child subservers can be returned in any order. For example:
// Given:
// A - B - C
// .   |
// .   D - E
// Where again B is running this function, all valid return values are:
// {B, D, C, E}
// {B, D, C, E}
// {B, D, E, C}
//
// Ultimately serves as a translation layer between internal server representation and actual tree structure
// func (st *ServerTree) Get(rootName string) []*ServerTreeItem {
// 	treeItems := make([]*ServerTreeItem, 0)

// 	parentMap := make(map[string][]*ServerTreeItem)
// 	parentMap[rootName] = nil

// 	for _, server := range st.Tree {
// 		if _, ok := parentMap[server.Parent]; !ok {
// 			parentMap[server.Parent] = make([]*ServerTreeItem, 0)
// 		}
// 		parentMap[server.Parent] = append(parentMap[server.Parent], server)
// 	}

// 	for _, svr := range st.Tree {

// 	}

// 	return treeItems
// }
