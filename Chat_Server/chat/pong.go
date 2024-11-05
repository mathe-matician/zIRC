package chat

import "github.com/phuslu/log"

func pong(params map[string]interface{}) Response {
	msg := "Running PONG..."
	log.Info().Msg(msg)

	_client, ok := params["client"]
	if _client == nil || !ok {
		log.Error().Msg("Client not passed to JOIN command!!")
		return ERR_UNKNOWNERROR("")
	}
	client := _client.(*Client)

	_cmd_params := params["params"]
	if _cmd_params == nil {
		log.Debug().Msgf("NO params passed to PONG")
		return ERR_NEEDMOREPARAMS("")
	}

	cmd_params := _cmd_params.(string)
	log.Info().Msgf("cmd_params: %s, len: %d", cmd_params, len(cmd_params))

	if len(cmd_params) == 0 {
		log.Debug().Msgf("No params passed to PONG")
		return ERR_NEEDMOREPARAMS("")
	}

	client.PingPongChan <- cmd_params

	return EMPTY_RESPONSE()
}
