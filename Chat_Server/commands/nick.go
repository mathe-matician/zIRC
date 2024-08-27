package commands

import (
	c "zirc/client"

	"github.com/phuslu/log"
)

func nick(params map[string]interface{}) *map[string]string {
	log.Info().Msg("Running NICK...")

	_client, ok := params["client"]
	if !ok {
		log.Error().Msg("Client not passed to PASS command!!")
		return ERR_UNKNOWNERROR("")
	}

	nick, ok := params["params"]
	if !ok {
		log.Error().Msg("Params not in map!")
		return ERR_NEEDMOREPARAMS("")
	}

	client := _client.(*c.Client)
	client.SetNick(nick.(string))

	// no response from NICK signals success
	return EMPTY_RESPONSE()
}
