package chat

// LIST [<channel>{,<channel>}] [<server>]
// channel — Optional. A comma-separated list of channels to query.
// server — Optional. Used in some IRC networks to specify which server to query.
// E.g.
//
//	LIST
//	LIST #chatroom1,#chatroom2
//
// The server responds with numeric replies:
// 321 — Start of list (header).
// 322 — A channel entry: includes channel name, user count, and topic.
//
//	:irc.example.com 322 zach #general 42 :General discussion
//
// 323 — End of list.
//
// If the server supports IRCv3 capabilities, it may:
// Attach message tags to 322 messages.
// Use LIST responses with additional metadata (e.g., via IRCv3 LIST-EXTENDED).
// e.g.
//
//	LIST >10
//	LIST m #python
func list(params map[string]interface{}) Response {
	// list all channels on this server
	// Forward LIST to all other servers

	// TODO
	// talk to ChanServ
	return EMPTY_RESPONSE()
}
