package client

import (
	"strconv"
	rc "zirc/remote_conn"

	"github.com/google/uuid"
	"github.com/phuslu/log"
)

type Session struct {
	id            uuid.UUID
	end_timestamp *uuid.Time
}

type Client struct {
	nick    string
	user    string
	session Session
	conn    *rc.RemoteConn
	send    chan string
}

func stringTimeFromUnixTimestamp(time uuid.Time) string {
	timestamp_s, timestamp_ns := time.UnixTime()
	timestamp := strconv.FormatInt(timestamp_s, 10) + "." + strconv.FormatInt(timestamp_ns, 10)
	return timestamp
}

func NewClient(nick string, user string, conn *rc.RemoteConn) (*Client, *string, error) {
	s, err := NewSession()
	if err != nil {
		log.Error().Msgf("Error creating new client %s", err.Error())
		return nil, nil, err
	}
	session_timestamp := stringTimeFromUnixTimestamp(s.id.Time())

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

	return &Session{
		id:            uuid,
		end_timestamp: nil,
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

func (c *Client) SetNick(nick string) {
	c.nick = nick
}

func (c *Client) SetUser(user string) {
	c.user = user
}
