package chat

import (
	"strings"

	"github.com/phuslu/log"
)

var cap_subcommands = map[string]string{
	"LS":    "", // any
	"LIST":  "", // any (same as LS)
	"REQ":   "", // any
	"END":   "", // any
	"CLEAR": "", // any
	"ACK":   "", // server only
	"NAK":   "", // server only
	"DEL":   "", // server only
	"NEW":   "", // server only
}

func cap(params map[string]interface{}) Response {
	msg := "Running CAP..."
	log.Info().Msg(msg)

	_params := params["params"]
	p := _params.(string)

	log.Debug().Msgf("CAP params: %s", p)

	if len(p) == 0 {
		return ERR_NEEDMOREPARAMS("")
	}

	p_break := strings.Split(p, " ")

	// trim off sub_cmd
	sub_cmd, p_break := p_break[0], p_break[1:]

	res := Reply{
		code: "*",
		msg:  "",
	}

	_client, ok := params["client"]
	if _client == nil || !ok {
		return ERR_UNKNOWNERROR("")
	}
	client := _client.(*Client)
	if client == nil {
		return ERR_UNKNOWNERROR("")
	}

	if sub_cmd == "LS" {
		log.Debug().Msgf("CAP LS")
		res.MsgOverride(cap_ls())
	} else if sub_cmd == "LIST" {
		log.Debug().Msgf("CAP LIST")
		res.MsgOverride(cap_list(&client.CapState))
	} else if sub_cmd == "REQ" {
		log.Debug().Msgf("CAP REQ")
		cap_req(p_break, &client.CapState)
	} else if sub_cmd == "ACK" {
		// only server for responses
		log.Debug().Msgf("CAP ACK")
	} else if sub_cmd == "NAK" {
		// only server for responses
		log.Debug().Msgf("CAP NAK")
	} else if sub_cmd == "CLEAR" {
		log.Debug().Msgf("CAP CLEAR")
		cap_clear(&client.CapState)
	} else if sub_cmd == "END" {
		log.Debug().Msgf("CAP END")
		cap_end(&client.CapState, &client.Capabilities)
	} else if sub_cmd == "NEW" {
		log.Debug().Msgf("CAP NEW")
	} else if sub_cmd == "DEL" {
		log.Debug().Msgf("CAP DEL")
	} else {
		// unknown CAP subcommand
		log.Debug().Msgf("Unknown CAP subcommand: %s", sub_cmd)
		return ERR_UNKNOWNCOMMAND("")
	}

	return &res
}

// Usage: CAP LS or CAP LS :302
// Description: This command is used by the client to request a list of capabilities supported by the server. The server responds with the available capabilities.
// The optional 302 parameter indicates that the client is compatible with IRCv3.2+ features.
func cap_ls() string {
	return "CAP * LS :" + (*g_Server).Config["CAPABILITIES"]
}

func cap_list(client_caps *string) string {
	log.Debug().Msgf("CAP LIST inner. client_caps: %s", *client_caps)
	return "CAP * LIST :" + *client_caps
}

// Usage: CAP REQ :<capability> [<capability>...]
// Description: The client uses this command to request a specific set of capabilities from the server. If the server supports the requested capabilities, it enables them for the client. Multiple capabilities can be requested at once, separated by spaces.
// CAP REQ :multi-prefix sasl
func cap_req(params []string, client_caps *string) {
	log.Debug().Msgf("CAP REQ inner. Params: %s, client_caps: %s", params, *client_caps)
	// respond with CAP ACK if all works out
	for _, c := range params {
		if !g_Server.HasCapability(c) {
			continue
		}

		if strings.Contains(*client_caps, c) {
			continue
		}
		log.Debug().Msgf("Adding %s", c)
		*client_caps += c + " "
	}
}

// Usage: CAP ACK :<capability> [<capability>...]
// Description: The server uses this command to acknowledge that the capabilities requested by the client (via CAP REQ) have been enabled. The client will receive this message after a successful CAP REQ.
// CAP ACK :multi-prefix sasl
func cap_ack() {
	// :server_name CAP * ACK :<new_capability>
}

// Usage: CAP NAK :<capability> [<capability>...]
// Description: The server uses this command to inform the client that some or all of the requested capabilities in a CAP REQ command could not be enabled. This indicates a failed capability request.
// CAP NAK :unsupported-capability
func cap_nak() {

}

// Usage: CAP CLEAR
// Description: This command is used by the client to clear (disable) all capabilities that are currently enabled. The server responds by disabling all capabilities for the client, essentially resetting the connection to a state where no capabilities are enabled.
// CAP CLEAR
// no direct response to the user - just clear caps
func cap_clear(client_caps *string) {
	*client_caps = ""
}

// Usage: CAP END
// Description: This command is used by the client to indicate that it has finished capability negotiation. After sending CAP END, the client moves on to the next steps of the connection (such as sending NICK and USER). If the client does not require any capabilities or has finished negotiating, it sends CAP END to complete the process.
// CAP END
func cap_end(requested_caps *string, client_caps *map[string]string) {
	if requested_caps != nil && *requested_caps != "" {
		caps := strings.Split((*requested_caps), " ")
		for _, c := range caps {
			(*client_caps)[c] = ""
		}
	}
}

// Usage: CAP NEW :<capability> [<capability>...]
// Description: The server uses this command to notify the client of new capabilities that have become available after the initial capability negotiation. This is useful when new features are dynamically enabled on the server while a client is connected.
// CAP NEW :message-tags account-tag
func cap_new() {

}

// Usage: CAP DEL :<capability> [<capability>...]
// Description: The server uses this command to inform the client that certain capabilities are no longer available. This may occur if a server disables certain capabilities or features after the client has connected.
// CAP DEL :sasl
func cap_del() {

}

// SERVER
// TODO
// this is a much larger change to do dynamically...
// possibly impossible depending on the capability
func addCapability() {

}

///// CAPABILITIES
// 1. multi-prefix
// Description: Allows the server to send multiple user prefixes (e.g., @ for ops, + for voiced) in the same message. This means clients can see all user modes in a channel without separate messages for each mode.
// 2. sasl
// Description: Stands for Simple Authentication and Security Layer. This capability enables secure authentication mechanisms, such as PLAIN or EXTERNAL, which are often used to authenticate users via services like NickServ or external authentication systems.
// 3. account-notify
// Description: Informs the client when a user’s account status changes, such as when they log in or out of their registered account. Useful for tracking account-based statuses in real time.
// 4. away-notify
// Description: Notifies the client when a user sets or unsets their "away" status. This avoids the need for the client to repeatedly check the user’s status.
// 5. cap-notify
// Description: Allows the server to notify the client dynamically of new capabilities that become available after the initial capability negotiation (after the connection has been established).
// 6. extended-join
// Description: When a user joins a channel, this capability makes the server include the user’s account name and away status (if available) in the JOIN message. This gives more information to the client upon user join events.
// 7. chghost
// Description: Allows the server to notify clients when another user’s host or IP address changes. This can be useful for dynamic DNS or when users are connecting from multiple networks.
// 8. message-tags
// Description: Enables clients and servers to add metadata to messages in the form of tags. Tags can include timestamp information, message IDs, and other contextual data, allowing richer communication protocols between clients and bots or services.
// 9. invite-notify
// Description: Notifies clients in a channel when someone is invited to join that channel. This lets users track invitation actions in real time.
// 10. labeled-response
// Description: Allows clients to attach a label to requests, enabling them to match responses with their corresponding requests, even if the responses arrive out of order. This can be particularly useful for bots or clients handling many simultaneous requests.
// 11. server-time
// Description: Adds an accurate timestamp to messages, showing the time at which the message was received by the server. This helps synchronize events across clients, even if there's a lag or delay in message transmission.
// 12. echo-message
// Description: When enabled, the server sends a copy of each message a client sends back to that client. This is helpful for confirming that the message was processed and delivered by the server, especially in scenarios where network reliability is a concern.
// 13. batch
// Description: Enables the server to group multiple related messages together into a single batch. This reduces the need for clients to process events individually, streamlining message processing (e.g., for join/part messages on multiple channels).
// 14. setname
// Description: Allows clients to change their real name (or "gecos") field after they have already connected, without having to disconnect and reconnect.
// 15. tls
// Description: Signals that the connection is encrypted with Transport Layer Security (TLS). This capability provides a way to confirm secure communication between the server and client.
// 16. sts
// Description: Stands for Strict Transport Security. This capability informs the client that it should use TLS (encrypted communication) when connecting to the server and provides the duration for which this policy should be applied.
// 17. draft/message-ids
// Description: Assigns unique message IDs to each message sent by the server, allowing clients to track and manage message delivery and avoid duplicate message handling in the case of retransmissions.
// 18. account-tag
// Description: Attaches the user’s account name to their messages, making it easier for clients to identify messages from authenticated users (e.g., for use in logging, moderation, or distinguishing between registered/unregistered users).
// 19. metadata
// Description: Provides a way for clients to set, delete, and receive arbitrary metadata for users or channels. This could include information like display names, custom user statuses, or other server-specific metadata.
// 20. message-tags
// Description: Allows the client to send and receive custom tags on IRC messages. These tags can contain additional metadata or information related to the message (e.g., timestamps, user info).
// 21. znc.in/self-message
// Description: Allows clients connected through a ZNC bouncer to see their own messages, even if they are buffered by the bouncer and delivered later.
// 22. draft/language
// Description: Enables users to specify the language they prefer for server messages. This can be useful in multilingual environments or services that provide localized responses.

func sasl() {

}
