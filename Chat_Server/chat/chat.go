package chat

import (
	"bytes"
	"errors"
	"strings"

	"zirc/commands"
	sm "zirc/servermanager"

	"github.com/phuslu/log"
)

func ProcessMessage(recv_buf *[]byte, server_manager *sm.ServerManager) []byte {
	log.Debug().Msg("------------MSG START------------")
	trimmed_msg := string(bytes.Trim(bytes.TrimLeft(*recv_buf, " "), "\x00"))
	log.Info().Msgf("Raw Client msg: %s", trimmed_msg)

	split_msg := strings.Split(trimmed_msg, " ")
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
		processTags(split_msg[0])
		// TODO - remove tags so the next chunk is the optional source
		split_msg = split_msg[1:]
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

	// TODO - process command as it should be the next thing in the message
	log.Debug().Msgf("Validating command %s", split_msg[0])
	cmd, err := commands.VerifyCommand(split_msg[0])
	if err != nil {
		return []byte(err.Error())
	}

	res := cmd.Fn()

	log.Debug().Msg("------------MSG END------------")
	// msg := "Server echo cmd: " + split_msg[0]
	b := []byte(res)
	return b
}
