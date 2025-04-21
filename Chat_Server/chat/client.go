package chat

import (
	"fmt"
	"net"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	rc "zirc/remote_conn"

	"github.com/google/uuid"
	"github.com/phuslu/log"
)

const MAX_NICK_LEN = 32

type StateItem struct {
	Value interface{}
	Res   Response
}

type Session struct {
	id            uuid.UUID
	end_timestamp *uuid.Time
	state         map[string]string
}

type AuthNType string
type AuthZType string

type ClientAuth struct {
	AuthenticationState map[string]int
	IsAuthenticated     bool
	AuthenticationType  AuthNType
	AuthorizationType   AuthZType
}

type Client struct {
	UID           uuid.UUID
	nick          string
	NickTimestamp time.Time
	user          string
	RealName      string
	Registered    bool
	server        string
	session       Session
	conn          *rc.RemoteConn
	ClientConn    *net.Conn
	send          chan string
	UserModes     []Mode
	Channels      []*Channel
	PrivateConvos []string // TODO - idk what this structure / process looks like / or even if it matters
	Host          string
	AwayMessage   string
	Capabilities  map[string]string
	CapState      string
	Auth          ClientAuth
	PingPongChan  chan string
	IsServer      bool
}

func NewClientAuth() ClientAuth {
	auth_state := make(map[string]int)
	return ClientAuth{
		AuthenticationState: auth_state,
		IsAuthenticated:     false,
		AuthenticationType:  "",
		AuthorizationType:   "",
	}
}

func NewStateItem(value interface{}) StateItem {
	return StateItem{
		Value: value,
	}
}

func stringTimeFromUnixTimestamp(time uuid.Time) string {
	timestamp_s, timestamp_ns := time.UnixTime()
	timestamp := strconv.FormatInt(timestamp_s, 10) + "." + strconv.FormatInt(timestamp_ns, 10)
	return timestamp
}

// NewClient creats a new client struct
// This is run on the current server, so IRC_SERVER_DNS_NAME
// will be set to the server's name
func NewClient(nick string, user string, conn *rc.RemoteConn, Conn *net.Conn, isServer bool) (*Client, *string, error) {
	s, err := NewSession()
	if err != nil {
		log.Error().Msgf("Error creating new client %s", err.Error())
		return nil, nil, err
	}
	session_timestamp := stringTimeFromUnixTimestamp(s.id.Time())

	s_chan := make(chan string)
	// by default add these user modes to all new clients
	// C: don't allow CTCP
	// i: invisible? hide them from all other users UNLESS they are on the same channel?
	// x: cloaked mode, cloak the ip address
	user_modes := []Mode{
		{ModeChar: "C", Params: ""},
		{ModeChar: "x", Params: ""},
	}

	caps := make(map[string]string)
	// authState := make(map[string]int)

	// maxUserChans := G_Config.Server.Max_user_channels

	// TODO
	// clients need to hold state of joined channels
	// for queries like who
	// and to keep track of if they can join more channels or not
	/// EDIT/
	// this already exists below - client.Channels

	uuid, err := uuid.NewV7()
	if err != nil {
		log.Error().Msgf("NewClient: Error generating uuid %s", err.Error())
		return nil, nil, err
	}

	serverName := G_Config.Server.Server_name
	if isServer {
		serverName = ""
	}

	return &Client{
		UID:          uuid,
		nick:         nick,
		user:         user,
		Registered:   false,
		server:       serverName,
		session:      *s,
		ClientConn:   Conn,
		conn:         conn,
		send:         s_chan,
		UserModes:    user_modes,
		Host:         serverName,
		AwayMessage:  "",
		Capabilities: caps,
		Channels:     make([]*Channel, 0),
		Auth:         NewClientAuth(),
		PingPongChan: make(chan string, 1), // buffered channel
		IsServer:     isServer,
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

func (c *Client) IsTarget() {}

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

func (c *Client) GetState(key string) string {
	if val, ok := c.session.state[key]; ok {
		return val
	}
	return ""
}

func (c *Client) AddState(key, value string) {
	c.session.state[string(key)] = string(value)
}

func (c *Client) LogState() {
	state := ""
	for key, value := range c.session.state {
		state += fmt.Sprintf("key: %s, value: %s\n", key, value)
	}
	log.Debug().EmbedObject(c).Msgf("STATE: %s", state)
}

func (c *Client) RemoveState(key string) {
	delete(c.session.state, key)
}

// target format :nickname!username@hostname
func (c *Client) FormattedClientDetails() string {
	ip := c.conn.Ip
	if c.HasUserMode("x") {
		ip = "cloak.z.irc"
	}
	return fmt.Sprintf(":%s!%s@%s", c.nick, c.user, ip)
}

func (c *Client) AddMode(mode Mode) {
	// c.mu.Lock()
	// defer c.mu.Unlock()
	char := mode.ModeChar
	if c.HasUserMode(char) {
		// if the user already has this mode don't add it again just do nothing
		return
	}

	c.UserModes = append(c.UserModes, mode)
}

func (c *Client) RemoveMode(mode string) {
	for i, m := range c.UserModes {
		if m.ModeChar == mode {
			c.UserModes = slices.Delete(c.UserModes, i, i+1)
		}
	}
}

func (c *Client) HasUserMode(mode string) bool {
	for _, m := range c.UserModes {
		if m.ModeChar == mode {
			return true
		}
	}
	return false
}

func (c *Client) HasCapability(cap string) bool {
	if _, ok := c.Capabilities[cap]; !ok {
		return false
	}
	return true
}

func (c *Client) GetConn() *rc.RemoteConn {
	return c.conn
}

func (c *Client) Nick() string {
	return c.nick
}

func (c *Client) Ip() string {
	if c.HasUserMode("x") {
		return "cloak.z.irc"
	}
	return c.conn.Ip
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

func (c *Client) FmtModes() string {
	modes := ""
	params := ""
	for _, m := range c.UserModes {
		modes += m.ModeChar
		params += " " + m.Params
	}
	return modes + " " + strings.TrimLeft(params, " ")
}

// ! @ # $ % & * + ( ) = / ? : ; , . < >
// Are invalid characters for nicks
// valid special {}[]-_^\
var nick_re = regexp.MustCompile(`^[A-Za-z_\-\[\]\\^\{\}][A-Za-z0-9_\-\[\]\\^\{\}]*$`)

func IsValidNickName(nick string) bool {
	if len(nick_re.FindAllStringSubmatch(nick, -1)) == 0 || len(nick) > MAX_NICK_LEN {
		return false
	}
	return true
}
