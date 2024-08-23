package chat

import (
	"bytes"

	"github.com/phuslu/log"
)

func ProcessMessage(recv_buf *[]byte) []byte {
	log.Debug().Msg("------------MSG START------------")
	trimmed_msg := string(bytes.Trim(bytes.TrimLeft(*recv_buf, " "), "\x00"))
	log.Info().Msgf("Raw Client msg: %s", trimmed_msg)

	log.Info().Msg("Before parsing tag data")
	// TODO - do tag parsing
	log.Debug().Msg("------------MSG END------------")
	b := []byte("hi from irc server")
	return b
}
