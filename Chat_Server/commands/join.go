package commands

import (
	"github.com/phuslu/log"
)

func join(params map[string]interface{}) *map[string]string {
	msg := "Running JOIN..."
	log.Info().Msg(msg)

	_task_runner := params["task_runner"]
	task_runner := _task_runner.(chan string)
	task_runner <- "RUN JOIN TASK!"

	response := map[string]string{
		"msg": msg,
	}
	return &response
}
