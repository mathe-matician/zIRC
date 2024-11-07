package chat

import (
	"fmt"
	"strconv"
	"time"
	"zirc/helpers"

	"github.com/google/uuid"
	"github.com/phuslu/log"
)

// TODO
// probably don't have to pass the entire client object here
// just the connection and the chan
// func ping_response_waiter(c *Client, ping_token string) {
// 	timeout, err := strconv.Atoi(helpers.GetEnv("IRC_PING_PONG_TIMEOUT", "10"))
// 	if err != nil {
// 		log.Error().Msgf("Error converting ping/pong timeout to int: %s", err)
// 		timeout = 15
// 	}

// 	for timeout != 0 {
// 		time.Sleep(1 * time.Second)
// 		timeout--

// 		pong_token := <-c.PingPongChan
// 		log.Debug().EmbedObject(c).Msgf("PONG token: %s", pong_token)
// 		if pong_token != ping_token {
// 			log.Debug().EmbedObject(c).Msgf("%s != %s", pong_token, ping_token)
// 			// if we get the wrong ping token then close the connection
// 			break
// 		}
// 	}

// 	log.Debug().EmbedObject(c).Msgf("Closing client connection")
// 	// If we get here the timeout was reached so close the client connection
// 	client_conn := c.ClientConn
// 	(*client_conn).Close()
// }

func ping(c *Client) {
	log.Debug().EmbedObject(c).Msgf("Running PING keepalive...")

	for {
		ping_token, err := uuid.NewV7()
		pingToken := ping_token.String()
		if err != nil {
			log.Error().EmbedObject(c).Msgf("Could not generate ping_token uuid!!!: %s", err.Error())
			pingToken = "TOKEN123"
		}

		conn := c.ClientConn
		cmdmsg := fmt.Sprintf("PING :%s", pingToken)

		task_runner := g_Server._MessageManager.Task_runner

		_task := NewTask(UNICAST, cmdmsg, 0.0, conn, nil, false)
		ping_task := []*Task{}
		ping_task = append(ping_task, _task)

		task_runner <- ping_task
		// TODO
		// add all parts of response (server, etc)
		// ping_msg := helpers.FormatResponse(cmdmsg)
		// if _, err := (*conn).Write([]byte(ping_msg)); err != nil {
		// 	log.Error().EmbedObject(c).Msgf("Error writing to client: %s", err.Error())
		// 	return
		// }

		timeout, err := strconv.Atoi(helpers.GetEnv("IRC_PING_PONG_TIMEOUT", "10"))
		if err != nil {
			timeout = 10
		}

		select {
		case pong_token := <-c.PingPongChan:
			pongToken := pong_token[1:]
			if pongToken != pingToken {
				log.Info().EmbedObject(c).Msgf("%s != %s", pongToken, pingToken)
				log.Info().EmbedObject(c).Msgf("%s len: %d", pongToken, len(pongToken))
				log.Info().EmbedObject(c).Msgf("%s len: %d", pingToken, len(pingToken))
				(*conn).Close()
			}
			log.Info().EmbedObject(c).Msgf("PONG successful")
		case <-time.After(time.Duration(timeout) * time.Second):
			log.Info().EmbedObject(c).Msgf("Server never received PONG, closing connection")
			(*conn).Close()
		}
	}
}
