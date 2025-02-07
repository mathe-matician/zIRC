package chat

import "github.com/phuslu/log"

func netinfo(params map[string]interface{}) Response {
	msg := "Running NETINFO..."
	log.Info().Msg(msg)
	res := Reply{
		code: "333",
		msg:  msg,
	}
	return &res
}
