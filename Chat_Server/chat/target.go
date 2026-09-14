package chat

import "net"

type Target interface {
	IsTarget()
	GetConn() net.Conn
	GetPingPongChan() chan string
}

type RemoteTask struct {
	DNS        string
	ClientNick string
	Msg        string
}

func (rt *RemoteTask) IsTarget()                    {}
func (rt *RemoteTask) GetConn() net.Conn            { return nil }
func (rt *RemoteTask) GetPingPongChan() chan string { return nil }
