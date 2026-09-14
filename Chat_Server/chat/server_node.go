package chat

import (
	"net"

	"github.com/phuslu/log"
)

// ServerNode is a "snapshot" view of a server with all the information needed
// for routing messages to remote servers
type ServerNode struct {
	Name              string
	SID               string
	Description       string
	DirectlyConnected bool
	IsMe              bool
	HopCount          int
	Parent            *ServerNode
	Conn              *ServerConn
	Servers           []*ServerNode
	PingPongChan      chan string // buffered channel
	Server            *IrcServer  // optional - could probably delete
}

func (sn *ServerNode) IsTarget() {}
func (sn *ServerNode) GetConn() net.Conn {
	return sn.Conn.Conn
}

func (sn *ServerNode) GetPingPongChan() chan string {
	return sn.PingPongChan
}

func (sn *ServerNode) MarshalObject(e *log.Entry) {
	e.Str("SID", sn.SID).Str("name", sn.Name).Str("description", sn.Description).Int("hopcount", sn.HopCount)
}

type ServerNodeOption func(*ServerNode)

func NewServerNode(opts ...ServerNodeOption) *ServerNode {
	log.Debug().Msg("NewServerNode: initializing opts...")
	sn := &ServerNode{}
	for _, opt := range opts {
		opt(sn)
	}
	return sn
}

func WithPingPongChan() ServerNodeOption {
	return func(sn *ServerNode) {
		sn.PingPongChan = make(chan string, 1) // buffered channel
	}
}

func WithSID(sid string) ServerNodeOption {
	return func(sn *ServerNode) { sn.SID = sid }
}

func WithName(name string) ServerNodeOption {
	return func(sn *ServerNode) { sn.Name = name }
}

func WithDescription(description string) ServerNodeOption {
	return func(sn *ServerNode) { sn.Description = description }
}

func WithDirectlyConnected(directlyConnected bool) ServerNodeOption {
	return func(sn *ServerNode) { sn.DirectlyConnected = directlyConnected }
}

func WithIsMe(isMe bool) ServerNodeOption {
	return func(sn *ServerNode) { sn.IsMe = isMe }
}

func WithHopCount(hopCount int) ServerNodeOption {
	return func(sn *ServerNode) { sn.HopCount = hopCount }
}

func WithParent(parent *ServerNode) ServerNodeOption {
	return func(sn *ServerNode) { sn.Parent = parent }
}

func WithConn(conn *ServerConn) ServerNodeOption {
	return func(sn *ServerNode) { sn.Conn = conn }
}

func WithServers(servers []*ServerNode) ServerNodeOption {
	return func(sn *ServerNode) {
		if servers == nil {
			servers = make([]*ServerNode, 0)
		}
		sn.Servers = servers
	}
}
