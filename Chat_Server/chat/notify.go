package chat

import "github.com/phuslu/log"

func notify(params map[string]interface{}) Response {
	msg := "Running NOTIFY..."
	log.Info().Msg(msg)
	res := Reply{
		code: "333",
		msg:  msg,
	}
	return &res
}
