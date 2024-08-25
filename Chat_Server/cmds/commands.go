package commands

import (
	"github.com/phuslu/log"
)

type Command struct{}

var command_list = []string{
	"ADMIN",
	"AWAY",
	"CAP",
}

func (c *Command) PRIVMSG() {
	log.Info().Msg("PRIVMSG start")
}

func (c *Command) NOTIFY() {

}
