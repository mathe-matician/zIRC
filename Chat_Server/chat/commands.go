package chat

import (
	"net"
	"sync"

	"zirc/helpers"

	"github.com/phuslu/log"
)

type CommandFunc func(map[string]interface{}) Response

type Command struct {
	Fn       CommandFunc
	Metadata map[string]string
	MustLock bool
	mu       sync.Mutex
}

var command_map = map[string]Command{
	"AUTHENTICATE": *NewCommand(authenticate, map[string]string{"cap_req": "sasl"}, true),
	"REGISTER":     *NewCommand(register, map[string]string{"cap_req": "account-registration"}, true),
	"VERIFY":       *NewCommand(verify, map[string]string{"cap_req": "account-registration"}, true),
	"CAP":          *NewCommand(cap, make(map[string]string), false),
	"ERROR":        *NewCommand(error_cmd, make(map[string]string), false),
	"NICK":         *NewCommand(nick, make(map[string]string), false),
	"PASS":         *NewCommand(pass, make(map[string]string), false),
	"PONG":         *NewCommand(pong, make(map[string]string), false),
	"JOIN":         *NewCommand(join, map[string]string{"auth_req": "true"}, true),
	"PROTOCTL":     *NewCommand(protoctl, make(map[string]string), false),
	"PRIVMSG":      *NewCommand(privmsg, map[string]string{"auth_req": "true"}, false),
	"MODE":         *NewCommand(mode, map[string]string{"auth_req": "true"}, true),
	"NOTIFY":       *NewCommand(notify, map[string]string{"auth_req": "true"}, false),
	"SERVER":       *NewCommand(server, make(map[string]string), false),
	"USER":         *NewCommand(user, make(map[string]string), false),
	"WHO":          *NewCommand(who, map[string]string{"auth_req": "true"}, false),
	"QUIT":         *NewCommand(quit, make(map[string]string), false),
	// "PING":         *NewCommand(ping, make(map[string]string), false),
	// "WEBIRC":       *NewCommand(webirc, make(map[string]string), false),
}

// TODO
// for certain commands to use the mutex,
// we need to lock the mutex within the function
// otherwise ALL commands will lock and unlock mutex which isn't needed
// e.g. sending chat messages via PRIVMSG doesn't need to lock the mutex and shouldn't!
// actually it might...
func (c *Command) Run(params map[string]interface{}) Response {
	if c.MustLock {
		c.mu.Lock()
		defer c.mu.Unlock()
	}
	return c.Fn(params)
}

func NewCommand(fn CommandFunc, metadata map[string]string, must_lock bool) *Command {
	return &Command{
		Fn:       fn,
		Metadata: metadata,
		MustLock: must_lock,
	}
}

func (c *Command) AddMetadata(cmd, key, value string) {
	command_map[cmd].Metadata[key] = value
}

func (c *Command) DeleteMetadata(cmd, key, value string) {
	delete(command_map[cmd].Metadata, value)
}

// CommandValidation focuses strictly on the command (not params) and does the follow:
// - ensures the command is a valid IRC command
// - checks whether the client is registered or not and limits commands based on that
// - checks whether the client is a server or a client and limits more commands based on that
func commandValidation(cmd, client_password_state string, client_registered bool, capabilities map[string]string) (*Command, Response) {
	server_password := G_Config.Server.Password_file
	if len(server_password) != 0 && cmd != "PASS" && client_password_state != "accepted" {
		return nil, ERR_PASSWDMISMATCH(":You need to send your password before registering")
	}

	return_cmd, ok := command_map[cmd]
	if !ok {
		return nil, ERR_UNKNOWNCOMMAND("")
	}

	// TODO - need to check if this connection is from a server here - is this a valid TODO anymore?

	_, auth_req := command_map[cmd].Metadata["auth_req"]
	if !client_registered && auth_req {
		log.Debug().Msgf("CommandValidation: ")
		return nil, ERR_NOTREGISTERED("")
	}

	_, has_sasl := capabilities["sasl"]
	if cmd == "AUTHENTICATE" && !has_sasl {
		return nil, ERR_NOSASL("")
	}

	log.Debug().Msgf("Valid command: %s", cmd)
	return &return_cmd, EMPTY_RESPONSE()
}

func WELCOME_WRAPPER(client_conn *net.Conn, server_name, server_version, server_creation_date, server_usermodes, server_channelmodes, client_nick, client_details string) []*Task {
	_001 := string(helpers.FormatResponse(server_name, "001", client_nick, RPL_WELCOME("", client_nick).Msg(), client_details))
	_002 := string(helpers.FormatResponse(server_name, "002", client_nick, RPL_YOURHOST("", server_name, server_version).Msg()))
	_003 := string(helpers.FormatResponse(server_name, "003", client_nick, RPL_CREATED("", server_creation_date).Msg()))
	_rpl_myinfo := RPL_MYINFO("", client_nick, server_name, server_version, server_usermodes, server_channelmodes).Msg()
	_004 := string(helpers.FormatResponse(server_name, "004", _rpl_myinfo))

	rpl_welcome := NewTask(UNICAST, _001, 0.0, client_conn, nil, false, "")
	rpl_yourhost := NewTask(UNICAST, _002, 0.0, client_conn, nil, false, "")
	rpl_created := NewTask(UNICAST, _003, 0.0, client_conn, nil, false, "")
	rpl_myinfo := NewTask(UNICAST, _004, 0.0, client_conn, nil, false, "")

	responses := []*Task{
		rpl_welcome,
		rpl_yourhost,
		rpl_created,
		rpl_myinfo,
	}

	return responses
}

// Although not commonly used, a client can send an ERROR message to notify the server of a fatal error condition.
func error_cmd(params map[string]interface{}) Response {
	msg := "Running ERROR..."
	log.Info().Msg(msg)
	// response := map[string]string{
	// 	"msg": msg,
	// }
	res := Reply{
		code: "333",
		msg:  msg,
	}
	return &res
}

// This command is used in some IRC networks to negotiate specific protocol features.
func protoctl(params map[string]interface{}) Response {
	msg := "Running PROTOCTL..."
	log.Info().Msg(msg)
	res := Reply{
		code: "333",
		msg:  msg,
	}
	return &res
}

// Used in some IRC networks to provide the client’s real IP address when connecting through a web proxy.
// func webirc(pararms []string) string {
// 	msg := "Running WEBIRC..."
// 	log.Info().Msg(msg)
// 	return msg
// }
