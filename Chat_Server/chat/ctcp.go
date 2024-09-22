package chat

import (
	"fmt"

	"github.com/phuslu/log"
)

// Client-to-client protocol (CTCP) constant
// any CTCP msg will start and end with this character
const ctcpDelimiter = "\x01"

var ctcp_map = map[string]Command{
	"VERSION":    *NewCommand(ctcp_version, map[string]string{"auth_req": "true"}, false),
	"PING":       *NewCommand(ctcp_ping, map[string]string{"auth_req": "true"}, false),
	"ACTION":     *NewCommand(ctcp_action, map[string]string{"auth_req": "true"}, false),
	"TIME":       *NewCommand(ctcp_time, map[string]string{"auth_req": "true"}, false),
	"FINGER":     *NewCommand(ctcp_finger, map[string]string{"auth_req": "true"}, false),
	"CLIENTINFO": *NewCommand(ctcp_clientinfo, map[string]string{"auth_req": "true"}, false),
	"SOURCE":     *NewCommand(ctcp_source, map[string]string{"auth_req": "true"}, false),
	"USERINFO":   *NewCommand(ctcp_userinfo, map[string]string{"auth_req": "true"}, false),
	"DCC":        *NewCommand(ctcp_dcc, map[string]string{"auth_req": "true"}, false),
	"ERRMSG":     *NewCommand(ctcp_errmsg, map[string]string{"auth_req": "true"}, false),
}

func parseCTCP() {
	log.Debug().Msgf("Parsing CTCP")
}

// Queries the client for version information.
// Usage: PRIVMSG <nick> :\x01VERSION\x01
// Response: NOTICE <nick> :\x01VERSION <client version>\x01
func ctcp_version(params map[string]interface{}) Response {
	return nil
}

// Measures latency between two clients by sending a timestamp and expecting it to be echoed back.
// Usage: PRIVMSG <nick> :\x01PING <timestamp>\x01
// Response: NOTICE <nick> :\x01PING <same timestamp>\x01
func ctcp_ping(params map[string]interface{}) Response {
	return nil
}

// Sends an action (often used with /me commands).
// Usage: PRIVMSG <channel/nick> :\x01ACTION <text>\x01
// Example: /me waves sends * username waves.
func ctcp_action(params map[string]interface{}) Response {
	return nil
}

// Requests the local time of the recipient.
// Usage: PRIVMSG <nick> :\x01TIME\x01
// Response: NOTICE <nick> :\x01TIME <time string>\x01
func ctcp_time(params map[string]interface{}) Response {
	return nil
}

// Asks for information about the user (name, idle time, etc.). Historically, this command was inspired by the Unix finger command.
// Usage: PRIVMSG <nick> :\x01FINGER\x01
// Response: NOTICE <nick> :\x01FINGER <user information>\x01
func ctcp_finger(params map[string]interface{}) Response {
	return nil
}

// Requests a list of supported CTCP commands or asks for details about a specific command.
// Usage: PRIVMSG <nick> :\x01CLIENTINFO\x01
// Response: NOTICE <nick> :\x01CLIENTINFO <list of supported commands>\x01
func ctcp_clientinfo(params map[string]interface{}) Response {
	return nil
}

// Requests the URL or location of the client's source code, typically in open-source projects.
// Usage: PRIVMSG <nick> :\x01SOURCE\x01
// Response: NOTICE <nick> :\x01SOURCE <URL>\x01
func ctcp_source(params map[string]interface{}) Response {
	return nil
}

// Requests information about the user, such as a personal description.
// Usage: PRIVMSG <nick> :\x01USERINFO\x01
// Response: NOTICE <nick> :\x01USERINFO <description>\x01
func ctcp_userinfo(params map[string]interface{}) Response {
	return nil
}

// Allows users to initiate a direct connection between clients for file transfers or chat. DCC is one of the more complex CTCP commands.
// Common subcommands include:
// DCC CHAT: Opens a direct text chat connection.
// DCC SEND: Initiates a file transfer.
// Usage: PRIVMSG <nick> :\x01DCC SEND <filename> <IP> <port> <filesize>\x01
func ctcp_dcc(params map[string]interface{}) Response {
	return nil
}

// Sends an error message in response to an invalid or unrecognized CTCP command.
// Usage: PRIVMSG <nick> :\x01ERRMSG <error message>\x01
// Response: NOTICE <nick> :\x01ERRMSG <error details>\x01
// this is actually sent from the client...
// but maybe we first check the client's msg here and then relay the message if it is correct
func ctcp_errmsg(params map[string]interface{}) Response {
	err_msg := "\x01ERRMSG ? :Unknown CTCP command\x01"
	_cmd := params["cmd"]
	if _cmd != nil {
		cmd := _cmd.(string)
		err_msg = fmt.Sprintf("\x01ERRMSG %s :Unknown CTCP command\x01", cmd)
	}

	res := ErrorResponse{
		code: "",
		msg:  err_msg,
	}
	return &res
}
