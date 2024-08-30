package commands

import "fmt"

// TODO - may use one of these
//
// type Response *map[string]string

type Response interface {
	MsgOverride(msg_override string)
	Code() string
	Msg() string
}

type ErrorResponse struct {
	code string
	msg  string
}

func (er *ErrorResponse) MsgOverride(msg_override string) {
	if len(msg_override) != 0 {
		er.msg = msg_override
	}
	er.msg = msg_override
}

func (er *ErrorResponse) Code() string {
	return er.code
}

func (er *ErrorResponse) Msg() string {
	return er.msg
}

type Reply struct {
	code string
	msg  string
}

func (r *Reply) MsgOverride(msg_override string) {
	if len(msg_override) != 0 {
		r.msg = msg_override
	}
}

func (r *Reply) Code() string {
	return r.code
}

func (r *Reply) Msg() string {
	return r.msg
}

func EMPTY_RESPONSE() Response {
	return &Reply{}
}

// hard coding CRLF since these messages aren't sent outside of registration confirmation
func RPL_WELCOME(msg_override, nick string) Response {
	return &Reply{
		code: "001",
		msg:  fmt.Sprintf(":Welcome to the IRC Network %s", nick),
	}
}

func RPL_YOURHOST(msg_override, server_name, server_version string) Response {
	return &Reply{
		code: "002",
		msg:  fmt.Sprintf(":Your host is %s, running version %s", server_name, server_version),
	}
}

func RPL_CREATED(msg_override, server_creation_date string) Response {
	return &Reply{
		code: "003",
		msg:  fmt.Sprintf(":This server was created %s", server_creation_date),
	}
}

// usermodes: The list of available user modes (e.g., o, i, w, s).
// channelmodes: The list of available channel modes (e.g., o, p, s, m, t).
func RPL_MYINFO(msg_override, nick, server_name, server_version, usermodes, channelmodes string) Response {
	return &Reply{
		code: "004",
		msg:  fmt.Sprintf("%s %s %s %s %s", nick, server_name, server_version, usermodes, channelmodes),
	}
}

func RPL_TOPIC(msg_override, topic string) Response {
	return &Reply{
		code: "332",
		msg:  fmt.Sprintf(":irc.example.com 332 <nickname> <channel> :%s", topic),
		// msg:  ":irc.example.com 332 <nickname> <channel> :<topic>",
	}
}

func RPL_NAMREPLY(msg_override string) Response {
	return &Reply{
		code: "353",
		msg:  ":irc.example.com 353 <nickname> = <channel> :@<nick1> +<nick2> <nick3>",
	}
}

func RPL_ENDOFNAMES(msg_override string) Response {
	return &Reply{
		code: "366",
		msg:  ":irc.example.com 366 <nickname> <channel> :End of /NAMES list.",
	}
}

func ERR_UNKNOWNERROR(msg_override string) Response {
	er := &ErrorResponse{
		code: "400",
		msg:  ":Unknown error occurred",
	}
	er.MsgOverride(msg_override)
	return er
}

func ERR_UNKNOWNCOMMAND(msg_override string) Response {
	er := &ErrorResponse{
		code: "421",
		msg:  ":Unknown command",
	}
	er.MsgOverride(msg_override)
	return er
}

func ERR_NOTREGISTERED(msg_override string) Response {
	er := &ErrorResponse{
		code: "451",
		msg:  ":You have not registered",
	}
	er.MsgOverride(msg_override)
	return er
}

func ERR_NEEDMOREPARAMS(msg_override string) Response {
	er := &ErrorResponse{
		code: "461",
		msg:  ":Need more params",
	}
	er.MsgOverride(msg_override)
	return er
}

func ERR_ALREADYREGISTRED(msg_override string) Response {
	er := &ErrorResponse{
		code: "462",
		msg:  ":You may not reregister",
	}
	er.MsgOverride(msg_override)
	return er
}

func ERR_YOUREBANNEDCREEP(msg_override string) Response {
	er := &ErrorResponse{
		code: "463",
		msg:  ":You are banned from this server",
	}
	er.MsgOverride(msg_override)
	return er
}

func ERR_PASSWDMISMATCH(msg_override string) Response {
	er := &ErrorResponse{
		code: "464",
		msg:  ":Password incorrect",
	}
	er.MsgOverride(msg_override)
	return er
}
