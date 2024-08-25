package remote_conn

import (
	"github.com/phuslu/log"
)

type RemoteConn struct {
	Host string
	Ip   string
	Port string
}

func (c *RemoteConn) MarshalObject(e *log.Entry) {
	e.Str("host", c.Host).Str("ip", c.Ip).Str("port", c.Port)
}

func NewRemoteConn(host, ip, port string) *RemoteConn {
	return &RemoteConn{
		Host: host,
		Ip:   ip,
		Port: port,
	}
}
