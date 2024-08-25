package client

import (
	"strconv"
	rc "zirc/remote_conn"

	"github.com/google/uuid"
	"github.com/phuslu/log"
)

type Session struct {
	id uuid.UUID
}

type Client struct {
	nick    string
	user    string
	session Session
	conn    *rc.RemoteConn
	send    chan string
}

func (c *Client) MarshalObject(e *log.Entry) {
	e.Str("nick", c.nick).Str("user", c.user).Str("session_id", c.session.id.String()).Str("host", c.conn.Host).Str("ip", c.conn.Ip).Str("port", c.conn.Port)
}

func NewClient(nick string, user string, conn *rc.RemoteConn) (*Client, *string, error) {
	s, err := NewSession()
	if err != nil {
		log.Error().Msgf("Error creating new client %s", err.Error())
		return nil, nil, err
	}
	timestamp_s, timestamp_ns := s.id.Time().UnixTime()
	session_timestamp := strconv.FormatInt(timestamp_s, 10) + "." + strconv.FormatInt(timestamp_ns, 10)

	// sessions := []Session{*s}

	s_chan := make(chan string)

	return &Client{
		nick:    nick,
		user:    user,
		session: *s,
		conn:    conn,
		send:    s_chan,
	}, &session_timestamp, nil
}

func NewSession() (*Session, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		log.Error().Msgf("Error generating uuid %s", err.Error())
		return nil, err
	}

	return &Session{id: uuid}, nil
}

func (c *Client) SetNick(nick string) {
	c.nick = nick
}

func (c *Client) SetUser(user string) {
	c.user = user
}
