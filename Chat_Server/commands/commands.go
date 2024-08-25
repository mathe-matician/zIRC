package commands

import (
	"errors"

	"github.com/phuslu/log"
)

type CommandFunc func() string

// type CommandFunc func()

type Command struct {
	Fn       CommandFunc
	Metadata map[string]string
}

var command_map = map[string]Command{
	"JOIN":    *NewCommand(join),
	"PRIVMSG": *NewCommand(privmsg),
	"NOTIFY":  *NewCommand(notify),
}

func NewCommand(fn CommandFunc) *Command {
	return &Command{
		Fn:       fn,
		Metadata: map[string]string{"name": "default"},
	}
}

func VerifyCommand(cmd string) (*Command, error) {
	if val, ok := command_map[cmd]; ok {
		log.Debug().Msgf("Valid command: %s", cmd)
		// return_cmd := val.(Command) // needs type asseration as the map value is "any"
		return_cmd := val
		return &return_cmd, nil
	}
	return nil, errors.New("unknown command")
}

func join() string {
	msg := "Running JOIN..."
	log.Info().Msg(msg)
	return msg
}

func privmsg() string {
	msg := "Running PRIVMSG..."
	log.Info().Msg(msg)
	return msg
}

func notify() string {
	msg := "Running NOTIFY..."
	log.Info().Msg(msg)
	return msg
}
