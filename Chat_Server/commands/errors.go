package commands

func error_map_generator(err_code, msg, msg_override string) *map[string]string {
	res := map[string]string{
		"err_code": err_code,
		"msg":      msg,
	}
	if len(msg_override) != 0 {
		res["msg"] = msg_override
	}
	return &res
}

func EMPTY_RESPONSE() *map[string]string {
	return error_map_generator("", "", "")
}

func ERR_UNKNOWNERROR(msg_override string) *map[string]string {
	return error_map_generator("400", ":Unknown error occurred", msg_override)
}

func ERR_UNKNOWNCOMMAND(msg_override string) *map[string]string {
	return error_map_generator("421", ":Unknown command", msg_override)
}

func ERR_NOTREGISTERED(msg_override string) *map[string]string {
	return error_map_generator("451", ":You have not registered", msg_override)
}

func ERR_NEEDMOREPARAMS(msg_override string) *map[string]string {
	return error_map_generator("461", ":Need more params", msg_override)
}

func ERR_ALREADYREGISTRED(msg_override string) *map[string]string {
	return error_map_generator("462", ":You may not reregister", msg_override)
}

func ERR_YOUREBANNEDCREEP(msg_override string) *map[string]string {
	return error_map_generator("463", ":You are banned from this server", msg_override)
}

func ERR_PASSWDMISMATCH(msg_override string) *map[string]string {
	return error_map_generator("464", ":Password incorrect", msg_override)
}
