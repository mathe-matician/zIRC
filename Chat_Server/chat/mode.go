package chat

import (
	"github.com/phuslu/log"
)

func mode(params map[string]interface{}) Response {
	log.Info().Msgf("MODE sttart")
	msg := "Running MODE..."
	log.Info().Msg(msg)

	_cmd_params := params["params"]
	if _cmd_params == nil {
		log.Debug().Msgf("no params passed to MODE")
		return ERR_NEEDMOREPARAMS("")
	}

	cmd_params := _cmd_params.(string)
	log.Info().Msgf("cmd_params: %s, len: %d", cmd_params, len(cmd_params))

	// _client, ok := params["client"]
	// if _client == nil || !ok {
	// 	log.Error().Msg("Client not passed to MODE command!!")
	// 	return ERR_UNKNOWNERROR("")
	// }
	// client := _client.(*c.Client)

	if len(cmd_params) == 0 {
		return ERR_NEEDMOREPARAMS("")
	}

	_task_runner := params["task_runner"]
	task_runner := _task_runner.(chan []*Task)

	task_runner <- []*Task{}

	return EMPTY_RESPONSE()
}
