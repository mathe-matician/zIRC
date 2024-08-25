package server

import (
	"zirc/remote_conn"
)

type Server struct {
	Conn *remote_conn.RemoteConn
}

func (s *Server) PRIVMSG() {

}
