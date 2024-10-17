package chat

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
	code           string
	msg            string
	show_client_ip bool
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
	code           string
	msg            string
	show_client_ip bool
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

func RPL_WELCOME(msg_override, nick string) Response {
	return &Reply{
		code:           "001",
		msg:            fmt.Sprintf(":Welcome to the IRC Network %s", nick),
		show_client_ip: true,
	}
}

func RPL_YOURHOST(msg_override, server_name, server_version string) Response {
	return &Reply{
		code:           "002",
		msg:            fmt.Sprintf(":Your host is %s, running version %s", server_name, server_version),
		show_client_ip: true,
	}
}

func RPL_CREATED(msg_override, server_creation_date string) Response {
	return &Reply{
		code:           "003",
		msg:            fmt.Sprintf(":This server was created %s", server_creation_date),
		show_client_ip: true,
	}
}

// usermodes: The list of available user modes (e.g., o, i, w, s).
// channelmodes: The list of available channel modes (e.g., o, p, s, m, t).
func RPL_MYINFO(msg_override, nick, server_name, server_version, usermodes, channelmodes string) Response {
	return &Reply{
		code:           "004",
		msg:            fmt.Sprintf("%s %s %s %s %s", nick, server_name, server_version, usermodes, channelmodes),
		show_client_ip: true,
	}
}

// response code, often used to list various server capabilities, including the channel modes that the server supports. The CHANMODES=b,k,l,imnpst part indicates the types of modes supported:
// :irc.example.com 005 <nickname> CHANMODES=b,k,l,imnpst CASEMAPPING=rfc1459 :are supported by this server
func RPL_ISUPPORT(msg_override, server_creation_date string) Response {
	return &Reply{
		code:           "005",
		msg:            fmt.Sprintf(":This server was created %s", server_creation_date),
		show_client_ip: true,
	}
}

// response code, listing the user modes currently set for the user (+iow might indicate invisible, operator, and wallops receipt modes).
// :irc.example.com 221 <nickname> :+iow
func RPL_UMODEIS(msg_override, nick, modes string) Response {
	return &Reply{
		code:           "221",
		msg:            fmt.Sprintf("%s +%s", nick, modes),
		show_client_ip: false,
	}
}

// Sent after server replies with all responses of 352
func RPL_ENDOFWHO(msg_override, nick, modes string) Response {
	return &Reply{
		code:           "315",
		msg:            "End of WHO list",
		show_client_ip: false,
	}
}

func RPL_CHANNELMODEIS(msg_override, channel, modes string) Response {
	return &Reply{
		code:           "324",
		msg:            fmt.Sprintf("%s +%s", channel, modes),
		show_client_ip: false,
	}
}

// is used to inform a user about the creation time of a specific channel. It is typically sent as part of the response to a /mode or /join command, alongside other information about the channel's modes.
// :irc.example.com 329 <nickname> #channel <timestamp>
// :irc.example.com 329 Alice #mychannel 1694018882
func RPL_CREATIONTIME(msg_override, channel, timestamp string) Response {
	return &Reply{
		code:           "329",
		msg:            fmt.Sprintf("%s %s", channel, timestamp),
		show_client_ip: true,
	}
}

func RPL_TOPIC(msg_override, topic string) Response {
	if len(topic) == 0 {
		topic = "No topic is set"
	}
	return &Reply{
		code: "332",
		msg:  fmt.Sprintf(":%s", topic),
		// msg:  ":irc.example.com 332 <nickname> <channel> :<topic>",
		show_client_ip: true,
	}
}

// :server: The server sending the reply.
// 352: The numeric code for the WHO reply.
// <requesting_user>: The nickname of the user who issued the WHO query.
// <channel>: The channel the user is in. If the query is not channel-specific, this could be *.
// <user>: The user’s username (ident).
// <host>: The user’s hostname or IP address.
// <server>: The name of the IRC server the user is connected to.
// <nick>: The nickname of the user.
// <H|G>: The user's "away" status. H means the user is not away ("Here"), and G means they are marked as away.
// [*]: If present, it indicates that the user is an IRC operator.
// [@|+]: If the user has channel operator status, this field is @; if they have voice privileges, it's +. If neither, this field is empty.
// :<hopcount>: The number of hops between the server and the user (used to represent network distance).
// <real name>: The "real name" (or GECOS) field from the user’s connection information.
// E.g. :irc.example.com 352 Bob #example alice alice.example.com irc.example.com Alice H@ :0 Alice Smith
func RPL_WHOREPLY(
	msg_override,
	requesting_user,
	channel,
	user,
	host,
	server,
	nick,
	away_status,
	is_operator,
	operator_status,
	hopcount,
	real_name string,
) Response {
	msg := fmt.Sprintf(":%s 352 %s %s %s %s %s %s %s %s %s :%s %s", server, requesting_user, channel, user, host, server, nick, away_status, is_operator, operator_status, hopcount, real_name)

	if len(msg_override) != 0 {
		msg = msg_override
	}

	return &Reply{
		code:           "352",
		msg:            msg,
		show_client_ip: true,
	}
}

func RPL_NAMREPLY(msg_override string) Response {
	return &Reply{
		code:           "353",
		msg:            ":irc.example.com 353 <nickname> = <channel> :@<nick1> +<nick2> <nick3>",
		show_client_ip: true,
	}
}

func RPL_ENDOFNAMES(msg_override string) Response {
	return &Reply{
		code:           "366",
		msg:            ":irc.example.com 366 <nickname> <channel> :End of /NAMES list.",
		show_client_ip: true,
	}
}

func ERR_UNKNOWNERROR(msg_override string) Response {
	er := &ErrorResponse{
		code:           "400",
		msg:            ":Unknown error occurred",
		show_client_ip: true,
	}
	if len(msg_override) != 0 {
		er.MsgOverride(msg_override)
	}
	return er
}

// Used with:
// PRIVMSG
// :irc.example.com 401 <nickname> #nonexistent :No such nick/channel
func ERR_NOSUCHNICK(msg_override string) Response {
	er := &ErrorResponse{
		code:           "401",
		msg:            ":No such nick",
		show_client_ip: true,
	}
	if len(msg_override) != 0 {
		er.MsgOverride(msg_override)
	}
	return er
}

// used with:
// MODE
// :irc.example.com 403 <nickname> #nonexistent :No such channel
func ERR_NOSUCHCHANNEL(msg_override, channel string) Response {
	er := &ErrorResponse{
		code:           "403",
		msg:            fmt.Sprintf("%s :No such channel", channel),
		show_client_ip: true,
	}
	if len(msg_override) != 0 {
		er.MsgOverride(msg_override)
	}
	return er
}

// Used with:
// PRIVMSG
// If an attempt is made to send the message to multiple users or channels in one command (which is not allowed), this error would be returned.
// No response for success: No response is given by the server for a successfully sent PRIVMSG.
// Errors: The server only sends a response if there is an issue, such as a non-existent user, lack of permissions, or other problems with the target.
func ERR_TOOMANYTARGETS(msg_override string) Response {
	er := &ErrorResponse{
		code:           "407",
		msg:            fmt.Sprintf(":Too many recipients."),
		show_client_ip: true,
	}
	if len(msg_override) != 0 {
		er.MsgOverride(msg_override)
	}
	return er
}

// Used with:
// PRIVMSG
// This error is sent if the PRIVMSG command is missing the message text (i.e., no message content is provided after the :).
func ERR_NOTEXTTOSEND(msg_override, target string) Response {
	er := &ErrorResponse{
		code:           "412",
		msg:            fmt.Sprintf("%s :No text to send", target),
		show_client_ip: true,
	}
	if len(msg_override) != 0 {
		er.MsgOverride(msg_override)
	}
	return er
}

// Used with:
// PRIVMSG
// This is returned if the message is addressed to an invalid hostname or domain.
func ERR_NOTOPLEVEL(msg_override string) Response {
	er := &ErrorResponse{
		code:           "413",
		msg:            fmt.Sprintf(":Invalid hostname or domain"),
		show_client_ip: true,
	}
	if len(msg_override) != 0 {
		er.MsgOverride(msg_override)
	}
	return er
}

// Used with:
// PRIVMSG
// This error occurs if the target of the message contains a wildcard in an invalid position, like trying to message *.com.
func ERR_WILDTOPLEVEL(msg_override string) Response {
	er := &ErrorResponse{
		code:           "414",
		msg:            fmt.Sprintf(":Invalid hostname or domain"),
		show_client_ip: true,
	}
	if len(msg_override) != 0 {
		er.MsgOverride(msg_override)
	}
	return er
}

func ERR_UNKNOWNCOMMAND(msg_override string) Response {
	er := &ErrorResponse{
		code:           "421",
		msg:            ":Unknown command",
		show_client_ip: true,
	}
	if len(msg_override) != 0 {
		er.MsgOverride(msg_override)
	}
	return er
}

func ERR_ERRONEUSNICKNAME(msg_override string) Response {
	er := &ErrorResponse{
		code:           "432",
		msg:            ":Erroneous nickname",
		show_client_ip: true,
	}
	if len(msg_override) != 0 {
		er.MsgOverride(msg_override)
	}
	return er
}

func ERR_NICKNAMEINUSE(msg_override string) Response {
	er := &ErrorResponse{
		code:           "433",
		msg:            ":Nickname is already in use",
		show_client_ip: true,
	}
	if len(msg_override) != 0 {
		er.MsgOverride(msg_override)
	}
	return er
}

// This error indicates a nickname collision, which can occur in scenarios where a user is trying to register a nickname that another user is also attempting to register simultaneously.
func ERR_NICKCOLLISION(msg_override string) Response {
	er := &ErrorResponse{
		code:           "436",
		msg:            ":Nickname is already in use",
		show_client_ip: true,
	}
	if len(msg_override) != 0 {
		er.MsgOverride(msg_override)
	}
	return er
}

// This error can occur if a user is changing their nickname too frequently. It prevents rapid nickname changes to avoid abuse.
func ERR_NICKTOOFAST(msg_override string) Response {
	er := &ErrorResponse{
		code:           "437",
		msg:            ":Nick change too fast",
		show_client_ip: true,
	}
	if len(msg_override) != 0 {
		er.MsgOverride(msg_override)
	}
	return er
}

// used with: MODE
func ERR_USERNOTINCHANNEL(msg_override string) Response {
	er := &ErrorResponse{
		code:           "441",
		msg:            ":They aren't on that channel",
		show_client_ip: true,
	}
	if len(msg_override) != 0 {
		er.MsgOverride(msg_override)
	}
	return er
}

func ERR_NOTREGISTERED(msg_override string) Response {
	er := &ErrorResponse{
		code:           "451",
		msg:            ":You have not registered",
		show_client_ip: true,
	}
	if len(msg_override) != 0 {
		er.MsgOverride(msg_override)
	}
	return er
}

func ERR_NEEDMOREPARAMS(msg_override string) Response {
	er := &ErrorResponse{
		code:           "461",
		msg:            ":Need more params",
		show_client_ip: true,
	}
	if len(msg_override) != 0 {
		er.MsgOverride(msg_override)
	}
	return er
}

func ERR_ALREADYREGISTRED(msg_override string) Response {
	er := &ErrorResponse{
		code:           "462",
		msg:            ":You may not reregister",
		show_client_ip: true,
	}
	if len(msg_override) != 0 {
		er.MsgOverride(msg_override)
	}
	return er
}

func ERR_YOUREBANNEDCREEP(msg_override string) Response {
	er := &ErrorResponse{
		code:           "463",
		msg:            ":You are banned from this server",
		show_client_ip: true,
	}
	if len(msg_override) != 0 {
		er.MsgOverride(msg_override)
	}
	return er
}

func ERR_PASSWDMISMATCH(msg_override string) Response {
	er := &ErrorResponse{
		code:           "464",
		msg:            ":Password incorrect",
		show_client_ip: true,
	}
	if len(msg_override) != 0 {
		er.MsgOverride(msg_override)
	}
	return er
}

// used with:
// MODE
func ERR_UNKNOWNMODE(msg_override, char string) Response {
	er := &ErrorResponse{
		code:           "472",
		msg:            fmt.Sprintf("%s :is unknown mode char to me", char),
		show_client_ip: true,
	}
	if len(msg_override) != 0 {
		er.MsgOverride(msg_override)
	}
	return er
}

// used with:
// MODE
// 476 <nickname> <channel> :Bad Channel Mask
func ERR_BADCHANMASK(msg_override string) Response {
	er := &ErrorResponse{
		code:           "476",
		msg:            ":Bad Channel Mask",
		show_client_ip: true,
	}
	if len(msg_override) != 0 {
		er.MsgOverride(msg_override)
	}
	return er
}

// used with:
// MODE
func ERR_NOTONCHANNEL(msg_override, channel string) Response {
	er := &ErrorResponse{
		code:           "489",
		msg:            fmt.Sprintf("%s :You're not channel operator", channel),
		show_client_ip: true,
	}
	if len(msg_override) != 0 {
		er.MsgOverride(msg_override)
	}
	return er
}

func ERR_SASLFAIL(msg_override string) Response {
	er := &ErrorResponse{
		code:           "904",
		msg:            ":SASL authentication not enabled",
		show_client_ip: true,
	}
	if len(msg_override) != 0 {
		er.MsgOverride(msg_override)
	}
	return er
}

// Used with:
// when client hasn't enabled sasl capability, e.g. CAP REQ sasl
func ERR_NOSASL(msg_override string) Response {
	er := &ErrorResponse{
		code:           "908",
		msg:            ":SASL authentication not enabled",
		show_client_ip: true,
	}
	if len(msg_override) != 0 {
		er.MsgOverride(msg_override)
	}
	return er
}

func RPL_SASLMECHS(msg_override string) Response {
	// TODO
	// grab existing sasl mechs and return them
	er := &ErrorResponse{
		code:           "909",
		msg:            ":SASL authentication not enabled",
		show_client_ip: true,
	}
	if len(msg_override) != 0 {
		er.MsgOverride(msg_override)
	}
	return er
}
