package commands

import (
	"errors"

	c "zirc/client"
	"zirc/helpers"

	"github.com/phuslu/log"
)

type CommandFunc func(map[string]interface{}) string

type Command struct {
	Fn       CommandFunc
	Metadata map[string]string
}

var command_map = map[string]Command{
	"AUTHENTICATE": *NewCommand(authenticate, map[string]string{"auth_req": "true"}),
	"CAP":          *NewCommand(cap, make(map[string]string)),
	"ERROR":        *NewCommand(error_cmd, make(map[string]string)),
	"PASS":         *NewCommand(pass, make(map[string]string)),
	"PING":         *NewCommand(ping, make(map[string]string)),
	"PONG":         *NewCommand(pong, make(map[string]string)),
	"JOIN":         *NewCommand(join, map[string]string{"auth_req": "true"}),
	"PROTOCTL":     *NewCommand(protoctl, make(map[string]string)),
	"PRIVMSG":      *NewCommand(privmsg, map[string]string{"auth_req": "true"}),
	"NOTIFY":       *NewCommand(notify, map[string]string{"auth_req": "true"}),
	"SERVER":       *NewCommand(server, make(map[string]string)),
	// "WEBIRC":       *NewCommand(webirc, make(map[string]string)),
	"QUIT": *NewCommand(quit, make(map[string]string)),
}

func NewCommand(fn CommandFunc, metadata map[string]string) *Command {
	return &Command{
		Fn:       fn,
		Metadata: metadata,
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
func CommandValidation(cmd string, client *c.Client) (*Command, error) {
	_client_password_state := client.GetState("server_password")
	client_password_state := ""
	if _client_password_state != nil {
		client_password_state = _client_password_state.(string)
	}

	server_password := helpers.GetEnv("IRC_SERVER_PASSWORD", "")
	if len(server_password) != 0 && cmd != "PASS" && client_password_state != "accepted" {
		return nil, errors.New(":You need to send your password before registering")
	}

	val, ok := command_map[cmd]
	if !ok {
		return nil, errors.New("unknown command")
	}

	// TODO - need to check if server here

	_, auth_req := command_map[cmd].Metadata["auth_req"]
	if (len(client.Nick()) == 0 || len(client.User()) == 0) && auth_req {
		return nil, errors.New("You have not registered")
	}

	log.Debug().Msgf("Valid command: %s", cmd)
	return_cmd := val
	return &return_cmd, nil
}

func authenticate(params map[string]interface{}) string {
	msg := "Running AUTHENTICATE..."
	log.Info().Msg(msg)
	return msg
}

func cap(params map[string]interface{}) string {
	msg := "Running CAP..."
	log.Info().Msg(msg)
	return msg
}

// Although not commonly used, a client can send an ERROR message to notify the server of a fatal error condition.
func error_cmd(params map[string]interface{}) string {
	msg := "Running ERROR..."
	log.Info().Msg(msg)
	return msg
}

func ping(params map[string]interface{}) string {
	msg := "Running PING..."
	log.Info().Msg(msg)
	return msg
}

func pong(params map[string]interface{}) string {
	msg := "Running PONG..."
	log.Info().Msg(msg)
	return msg
}

func join(params map[string]interface{}) string {
	msg := "Running JOIN..."
	log.Info().Msg(msg)
	return msg
}

func privmsg(params map[string]interface{}) string {
	msg := "Running PRIVMSG..."
	log.Info().Msg(msg)
	return msg
}

// This command is used in some IRC networks to negotiate specific protocol features.
func protoctl(params map[string]interface{}) string {
	msg := "Running PROTOCTL..."
	log.Info().Msg(msg)
	return msg
}

func notify(params map[string]interface{}) string {
	msg := "Running NOTIFY..."
	log.Info().Msg(msg)
	return msg
}

// In the case of a server connection, this command can be used for server-to-server communications (typically not used by clients).
func server(params map[string]interface{}) string {
	msg := "Running SERVER..."
	log.Info().Msg(msg)
	return msg
}

// Used in some IRC networks to provide the client’s real IP address when connecting through a web proxy.
// func webirc(pararms []string) string {
// 	msg := "Running WEBIRC..."
// 	log.Info().Msg(msg)
// 	return msg
// }

func quit(params map[string]interface{}) string {
	msg := "Running QUIT..."
	log.Info().Msg(msg)
	return msg
}
