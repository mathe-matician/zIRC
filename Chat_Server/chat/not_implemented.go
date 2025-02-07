package chat

import "github.com/phuslu/log"

// In the case of a server connection, this command can be used for server-to-server communications (typically not used by clients).
func not_implemented(params map[string]interface{}) Response {
	msg := "This functionality is not imtplemented yet"
	log.Info().Msg(msg)
	res := ErrorResponse{
		code: "-666",
		msg:  msg,
	}
	return &res
}
