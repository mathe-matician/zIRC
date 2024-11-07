package chat

import (
	"net"

	"github.com/google/uuid"
	"github.com/phuslu/log"
)

const (
	SERVER        = "server"
	UNICAST       = "unicast"
	MULTICAST     = "multicast"
	BROADCAST     = "broadcast"
	SERVER_ACTION = "server_action"
)

// Locations can be Source, Destination, Channels, etc
type Location interface {
	NewLocation()
}

// type Channel struct {
// 	Prefix string
// 	Name   string
// }

func (c *Channel) NewLocation() {}

type Src struct {
	Conn *net.Conn
}

func (s *Src) NewLocation() {}

type Dest struct {
	Conn *net.Conn
}

func (d *Dest) NewLocation() {}

type Task struct {
	Id                  uuid.UUID
	Type                string
	Weight              float64
	Task                string
	ClientConn          *net.Conn
	Target              Target
	MultiCastSendToSelf bool
	// Src        Location
	// Dest       Location
}

func (t *Task) MarshalObject(e *log.Entry) {
	e.Str("id", t.Id.String()).Str("type", t.Type).Float64("weight", t.Weight).Str("task", t.Task)
}

// func NewChannel(prefix, name string) *Channel {
// 	return &Channel{
// 		Prefix: prefix,
// 		Name:   name,
// 	}
// }

func NewSrc(conn *net.Conn) *Src {
	return &Src{
		conn,
	}
}

func NewDest(conn *net.Conn) *Dest {
	return &Dest{
		conn,
	}
}

// func NewTask(_type, task string, src Location, dest Location, weight float64) *Task {
// NewTask[T int64 | float64](_type, task string, weight float64, client_conn *net.Conn, target *T)
func NewTask(_type, task string, weight float64, client_conn *net.Conn, target Target, mSendToSelf bool) *Task {
	uid, err := uuid.NewV7()
	if err != nil {
		return nil
	}

	return &Task{
		Id:                  uid,
		Type:                _type,
		Weight:              weight, // TODO - determine weight of task somehow, may be useful
		Task:                task,
		ClientConn:          client_conn,
		Target:              target,
		MultiCastSendToSelf: mSendToSelf,
	}
}
