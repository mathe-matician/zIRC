package commands

import (
	c "zirc/client"

	"github.com/phuslu/log"
)

func user(params map[string]interface{}) *map[string]string {
	log.Info().Msg("Running USER...")

	_client, ok := params["client"]
	if !ok {
		log.Error().Msg("Client not passed to USER command!!")
		return ERR_UNKNOWNERROR("")
	}

	user, ok := params["params"]
	if !ok {
		log.Error().Msg("Params not in map!")
		return ERR_NEEDMOREPARAMS("")
	}

	client := _client.(*c.Client)
	client.SetUser(user.(string))

	return EMPTY_RESPONSE()
}
