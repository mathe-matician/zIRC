package client

import (
	"strconv"

	"zirc/helpers"
	rc "zirc/remote_conn"

	"github.com/google/uuid"
	"github.com/phuslu/log"
)

type Session struct {
	id            uuid.UUID
	end_timestamp *uuid.Time
	state         map[string]string
}

type Client struct {
	nick    string
	user    string
	server  string
	session Session
	conn    *rc.RemoteConn
	send    chan string
}

func stringTimeFromUnixTimestamp(time uuid.Time) string {
	timestamp_s, timestamp_ns := time.UnixTime()
	timestamp := strconv.FormatInt(timestamp_s, 10) + "." + strconv.FormatInt(timestamp_ns, 10)
	return timestamp
}

// NewClient creats a new client struct
// This is run on the current server, so IRC_SERVER_DNS_NAME
// will be set to the server's name
func NewClient(nick string, user string, conn *rc.RemoteConn) (*Client, *string, error) {
	s, err := NewSession()
	if err != nil {
		log.Error().Msgf("Error creating new client %s", err.Error())
		return nil, nil, err
	}
	session_timestamp := stringTimeFromUnixTimestamp(s.id.Time())

	s_chan := make(chan string)

	return &Client{
		nick:    nick,
		user:    user,
		server:  helpers.GetEnv("IRC_SERVER_DNS_NAME", "localhost"),
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

	return &Session{
		id:            uuid,
		end_timestamp: nil,
		state:         make(map[string]string),
	}, nil
}

func (c *Client) MarshalObject(e *log.Entry) {
	e.Str("nick", c.nick).Str("user", c.user).Str("session_id", c.session.id.String()).Str("host", c.conn.Host).Str("ip", c.conn.Ip).Str("port", c.conn.Port)
}

func (c *Client) SetSessionEndTimestamp() (*string, error) {
	time, _, err := uuid.GetTime()
	if err != nil {
		log.Error().Msgf("Error generating uuid %s", err.Error())
		return nil, err
	}
	c.session.end_timestamp = &time
	end_timestamp := stringTimeFromUnixTimestamp(time)
	return &end_timestamp, nil
}

func (c *Client) GetState(key string) interface{} {
	if val, ok := c.session.state[key]; ok {
		return val
	}
	return nil
}

func (c *Client) UpdateState(key, value string) {
	c.session.state[string(key)] = string(value)
}

func (c *Client) Nick() string {
	return c.nick
}

func (c *Client) SessionId() string {
	return c.session.id.String()
}

func (c *Client) User() string {
	return c.user
}

func (c *Client) SetNick(nick string) {
	c.nick = nick
}

func (c *Client) SetUser(user string) {
	c.user = user
}
