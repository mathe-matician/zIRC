package chat

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"

	c "zirc/client"
	"zirc/commands"
	"zirc/helpers"
	sm "zirc/servermanager"

	"github.com/phuslu/log"
)

var re = regexp.MustCompile(`^\S*`) // captures until first space

func ProcessMessage(recv_buf *[]byte, client *c.Client, server_manager *sm.ServerManager) []byte {
	log.Debug().Msg("------------MSG START------------")
	trimmed_msg := string(bytes.Trim(bytes.TrimLeft(*recv_buf, " "), "\x00"))
	log.Info().Msgf("Raw Client msg: %s", trimmed_msg)

	split_msg := re.FindAllStringSubmatch(trimmed_msg, -1)[0]
	log.Debug().Msgf("Split Msg: %s", split_msg)
	if len(split_msg) == 0 {
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
		split_msg = split_msg[1:]
		split_msg = re.FindAllStringSubmatch(split_msg[0], -1)[0]
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
		log.Info().Msg("Message contains source prefix. Must be from another server...")
		split_msg = split_msg[1:]
		split_msg = re.FindAllStringSubmatch(split_msg[0], -1)[0]
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
	cmd, err := commands.CommandValidation(strings.Trim(split_msg[0], " "), client)
	if err != nil {
		return []byte(err.Error())
	}

	// remove command
	cmd_params := split_msg[1:]
	log.Debug().Msgf("Cmd params %s, len: ", cmd_params, len(cmd_params))

	// TODO - params aren't being passed correctly - something with REGEX

	// TODO - handle "chunked" messages where no CRLF exists - need to wait for the rest of the message
	//		timeout if the rest of the message doesn't come through - i.e. we don't get a CRLF in x seconds
	//
	// 		Also check that the client hasn't exceeded 8192 bytes (8 KB) which is the max message size

	res := cmd.Fn(cmd_params)

	log.Debug().Msg("------------MSG END------------")

	response := []byte(fmt.Sprintf(":%s code target %s \r\n", helpers.GetEnv("IRC_SERVER_DNS_NAME", "localhost"), res))
	return response
}
