package commands

import (
	t "zirc/task"

	"github.com/phuslu/log"
)

func join(params map[string]interface{}) Response {
	log.Info().Msgf("Join sttart")
	msg := "Running JOIN..."
	log.Info().Msg(msg)

	_cmd_params := params["params"]
	if _cmd_params == nil {
		log.Debug().Msgf("no params passed to join")
		return ERR_NEEDMOREPARAMS("")
	}

	cmd_params := _cmd_params.(string)
	log.Info().Msgf("cmd_params: %s, len: %d", cmd_params, len(cmd_params))

	if len(cmd_params) != 0 {
		_task_runner := params["task_runner"]
		task_runner := _task_runner.(chan []*t.Task)
		task := t.NewTask(t.MULTICAST, "RUN JOIN TASK!", 0.0)

		task_runner <- []*t.Task{task}
	}

	return EMPTY_RESPONSE()
}
