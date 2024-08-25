package client

import (
	rc "zirc/remote_conn"
)

type Client struct {
	Nick string
	User string
	Conn *rc.RemoteConn
	Send chan string
}
