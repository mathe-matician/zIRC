package chat

import "github.com/phuslu/log"

func quit(params map[string]interface{}) Response {
	msg := "Running QUIT..."
	log.Info().Msg(msg)
	res := Reply{
		code: "333",
		msg:  msg,
	}
	return &res
}
