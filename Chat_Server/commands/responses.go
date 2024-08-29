package commands

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

func (er ErrorResponse) MsgOverride(msg_override string) {
	if len(msg_override) != 0 {
		er.msg = msg_override
	}
}

func (er ErrorResponse) Code() string {
	return er.code
}

func (er ErrorResponse) Msg() string {
	return er.msg
}

type Reply struct {
	code string
	msg  string
}

func (r Reply) MsgOverride(msg_override string) {
	if len(msg_override) != 0 {
		r.msg = msg_override
	}
}

func (r Reply) Code() string {
	return r.code
}

func (r Reply) Msg() string {
	return r.msg
}

func EMPTY_RESPONSE() Response {
	return Reply{}
}

func RPL_WELCOME() Response {
	return Reply{}
}

func ERR_UNKNOWNERROR(msg_override string) Response {
	er := ErrorResponse{
		code: "400",
		msg:  ":Unknown error occurred",
	}
	er.MsgOverride(msg_override)
	return er
}

func ERR_UNKNOWNCOMMAND(msg_override string) Response {
	er := ErrorResponse{
		code: "421",
		msg:  ":Unknown command",
	}
	er.MsgOverride(msg_override)
	return er
}

func ERR_NOTREGISTERED(msg_override string) Response {
	er := ErrorResponse{
		code: "451",
		msg:  ":You have not registered",
	}
	er.MsgOverride(msg_override)
	return er
}

func ERR_NEEDMOREPARAMS(msg_override string) Response {
	er := ErrorResponse{
		code: "461",
		msg:  ":Need more params",
	}
	er.MsgOverride(msg_override)
	return er
}

func ERR_ALREADYREGISTRED(msg_override string) Response {
	er := ErrorResponse{
		code: "462",
		msg:  ":You may not reregister",
	}
	er.MsgOverride(msg_override)
	return er
}

func ERR_YOUREBANNEDCREEP(msg_override string) Response {
	er := ErrorResponse{
		code: "463",
		msg:  ":You are banned from this server",
	}
	er.MsgOverride(msg_override)
	return er
}

func ERR_PASSWDMISMATCH(msg_override string) Response {
	er := ErrorResponse{
		code: "464",
		msg:  ":Password incorrect",
	}
	er.MsgOverride(msg_override)
	return er
}
