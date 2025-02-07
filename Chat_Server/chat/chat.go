package chat

import (
	"errors"
	"reflect"
	"regexp"
	"strings"

	"zirc/helpers"

	"github.com/phuslu/log"
)

// captures until first space
var re = regexp.MustCompile(`^\S*`)

// captures until first space and has two capture groups, what comes before the space and what comes after
// e.g. `#general hello world` would contain two capture groups: ((#general), (hello world))
// meant to be used with FindStringSubmatch(str)
var cmd_re = regexp.MustCompile(`^(\S+)(.*)`)

// TODO
// need to pass a client OR server here
// it needs to be a Target
func ProcessMessage(trimmed_msg string, client *Client) []byte {
	log.Debug().Msg("------------MSG START------------")
	// trimmed_msg := string(bytes.Trim(bytes.TrimLeft(*recv_buf, " "), "\x00"))
	// log.Info().Msgf("Raw Client msg: %s", trimmed_msg)

	server := G_Config.Server.Dns_name
	target := "*"

	split_msg := re.FindAllStringSubmatch(trimmed_msg, -1)[0]
	log.Debug().Msgf("Split Msg: %s, len: %d", split_msg, len(split_msg))
	log.Debug().Msgf("Split Msg[0]: %s, len[0]: %d", split_msg[0], len(split_msg[0]))
	if len(split_msg) == 0 || len(split_msg[0]) == 0 {
		err := errors.New("message is empty")
		log.Error().Msg(err.Error())
		return []byte(err.Error())
	}

	// message format should look like the following:
	// message         ::= ['@' <tags> SPACE] [':' <source> SPACE] <command> <parameters> <crlf>
	//   SPACE           ::=  %x20 *( %x20 )   ; space character(s)
	//   crlf            ::=  %x0D %x0A        ; "carriage return" "linefeed"

	// parse tag data
	if string(split_msg[0][0]) == "@" {
		log.Info().Msg("Message has tag data. Processing source first.")
		processTags(strings.Trim(split_msg[0], " "))
		// TODO - remove tags so the next chunk is the optional source
		trimmed_msg = trimmed_msg[1:]
		split_msg = re.FindAllStringSubmatch(trimmed_msg, -1)[0]
	}

	if len(split_msg) == 0 {
		err := errors.New("message is empty")
		log.Error().Msg(err.Error())
		return []byte(err.Error())
	}

	// ONLY IF FROM ANOTHER SERVER
	// process optional source
	// only optional because if another IRC server is in this network
	// messages will be relayed from that server to this one
	// the message will have a prefix from that particular server
	// so we must check whether that is a valid server
	if string(split_msg[0][0]) == ":" {
		// Clients MUST NOT include a source when sending a message.
		// E.g. clients must be able to process messages whether from a server or client
		// TODO
		// check to see if the prefix is valid! if it isn't from a valid connected link
		// in the network, then it is a probably spoofed prefix
		// e.g. check the routing table
		log.Info().Msg("Message contains source prefix. Must be from another server...")
		log.Info().Msgf("Trimmed msg 1: %v", trimmed_msg)

		split_server_prefix := strings.Split(split_msg[0], ":")
		if len(split_server_prefix[1]) == 0 {
			// no server was included in the prefix - this should never happen
			err := errors.New("no server prefix included in message")
			log.Error().Msg(err.Error())
			return []byte(err.Error())
		}

		_, _, err := g_Server.RoutingTable.GetServer(split_server_prefix[1])
		if err != nil {
			// if the server doesn't exist in the routing table, this isn't a valid server
			// TODO
			// is there a race condition here?
			(*client.ClientConn).Close()
			return []byte("")
		}

		trimmed_msg = trimmed_msg[1:]
		split_msg = re.FindAllStringSubmatch(trimmed_msg, -1)[0]
	}

	if len(split_msg) == 0 {
		err := errors.New("message is empty")
		log.Error().Msg(err.Error())
		return []byte(err.Error())
	}

	// TODO - once client NICK / USER confirmed
	//		  synchronize with other IRC servers
	//		  notifying about the new connection / user information
	//		  via a broadcast or something
	//		  ---
	//		  So this server needs to be able to make updates to its own state mapping if necessary
	//		  from other servers
	//		  this should be handled in its own go routine most likely
	//		  e.g. user information, channel information
	//		  ---
	//		  S2S communication uses cmds like PING/PONG, SYNCHRONIZE

	log.Debug().Msgf("Validating command %s", split_msg[0])
	var cmd *Command
	var _validation_res Response

	strCmd := strings.Trim(split_msg[0], " ")

	if client.IsServer {
		cmd, _validation_res = serverCommandValidation(strCmd)
	} else {
		cmd, _validation_res = commandValidation(strCmd, client.session.state["server_password"], client.Registered, client.Capabilities)
	}

	if reflect.TypeOf(_validation_res).Name() == "ErrorResponse" || cmd == nil {
		log.Error().Msgf("Error during command validation")
		return helpers.FormatResponse(server, _validation_res.Code(), target, _validation_res.Msg())
	}

	// remove command
	str_cmd := re.FindAllStringSubmatch(trimmed_msg, -1)[0]
	log.Debug().Msgf("Cmd: %s, Cmd len: %d", str_cmd, len(str_cmd[0]))
	log.Debug().Msgf("trimmed_msg: %s, trimmed_msg len: %d", trimmed_msg, len(trimmed_msg))
	cmd_params := trimmed_msg[len(str_cmd[0])+1:]
	log.Debug().Msgf("Trimmed_msg AFTER regex && trim: %s", cmd_params)

	// TODO - params aren't being passed correctly - something with REGEX

	if !strings.HasSuffix(cmd_params, "\r\n") {
		// TODO - handle "chunked" messages where no CRLF exists - need to wait for the rest of the message
		//		timeout if the rest of the message doesn't come through - i.e. we don't get a CRLF in x seconds
		//
		// 		Also check that the client hasn't exceeded 8192 bytes (8 KB) which is the max message size
		// note we already define the max buffer in main.go - but double check the actualy max size 4k or 8k?
		log.Info().Msgf("TODO: Message has no CRLF... wait for rest of message!")
	}

	cmd_params = strings.TrimSuffix(cmd_params, "\r\n")

	cmd_param_slice := map[string]interface{}{}
	log.Debug().Msgf("Cmd params %s, len: %d", cmd_params, len(cmd_params))
	if len(cmd_params) != 0 {
		cmd_param_slice["params"] = strings.Trim(cmd_params, " ")
	}

	cmd_param_slice["client"] = client

	// log.Debug().Msgf("Before running CMD func")
	_response := cmd.Fn(cmd_param_slice)

	log.Info().Msgf("Command res msg: %s, code: %s", _response.Msg(), _response.Code())
	log.Debug().Msg("------------MSG END------------")

	msg := _response.Msg()
	code := _response.Code()

	if len(msg) == 0 {
		return []byte("")
	}

	if reflect.TypeOf(_response).Name() == "ErrorResponse" {
		return helpers.FormatResponse(server, code, msg, target)
	}

	// check for _response["target"] as some responses don't format target the same way
	if client.Registered {
		target = client.FormattedClientDetails()
		log.Debug().Msgf("chat.go target: %s", target)
	}
	// return helpers.FormatResponse(server, str_cmd[0], target, msg)
	finalMsg := helpers.FormatResponse(server, code, msg, target)
	log.Debug().Msgf("chat.go finalMsg: %s", finalMsg)
	return finalMsg
}
