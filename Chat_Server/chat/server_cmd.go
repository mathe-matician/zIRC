package chat

import "github.com/phuslu/log"

// In the case of a server connection, this command can be used for server-to-server communications (typically not used by clients).
func server(params map[string]interface{}) Response {
	msg := "Running SERVER..."
	log.Info().Msg(msg)
	res := Reply{
		code: "333",
		msg:  msg,
	}
	return &res
}
