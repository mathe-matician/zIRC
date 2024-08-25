package chat

import "github.com/phuslu/log"

func processTags(tags string) {
	log.Info().Msgf("Processing tags: %s", tags)

	// TODO - since tags are handled by the client
	//        we can check here whether they include the tags we are expecting.
}
