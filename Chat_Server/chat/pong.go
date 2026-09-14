package chat

import (
	"github.com/phuslu/log"
)

func pong(params map[string]interface{}) Response {
	// msg := "Running PONG..."
	// log.Info().Msg(msg)

	_client, ok := params["client"]
	if _client == nil || !ok {
		log.Error().Msg("Client not passed to PONG command!!")
		return ERR_UNKNOWNERROR("")
	}

	_cmd_params := params["params"]
	if _cmd_params == nil {
		log.Debug().Msgf("NO params passed to PONG")
		return ERR_NEEDMOREPARAMS("")
	}

	cmd_params := _cmd_params.(string)
	// log.Debug().Msgf("cmd_params: %s, len: %d", cmd_params, len(cmd_params))

	if len(cmd_params) == 0 {
		log.Debug().Msgf("No params passed to PONG")
		return ERR_NEEDMOREPARAMS("")
	}

	switch c := _client.(type) {
	case *Client:
		client := _client.(*Client)
		client.PingPongChan <- cmd_params
	case *ServerNode:
		server := _client.(*ServerNode)
		server.PingPongChan <- cmd_params
	default:
		log.Error().Msgf("Not a valid pong type: %T", c)
		return ERR_UNKNOWNERROR("")
	}

	return EMPTY_RESPONSE()
}
