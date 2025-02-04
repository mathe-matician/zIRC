package chat

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/phuslu/log"
)

func ping(c *Client) {
	// log.Debug().EmbedObject(c).Msgf("Running PING keepalive...")

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

		_task := NewTask(UNICAST, cmdmsg, 0.0, conn, nil, false, "")
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

		timeout := G_Config.Server.Ping_pong_timeout

		_duration := G_Config.Server.Ping_pong_timeout_duration
		var duration time.Duration
		if _duration == "minute" {
			duration = time.Minute
		} else if _duration == "second" {
			duration = time.Second
		}

		select {
		case pong_token := <-c.PingPongChan:
			pongToken := pong_token[1:]
			if pongToken != pingToken {
				log.Warn().EmbedObject(c).Msgf("%s != %s", pongToken, pingToken)
				(*conn).Close()
			}
			// log.Debug().EmbedObject(c).Msgf("PONG successful")
		case <-time.After(time.Duration(timeout) * duration):
			log.Info().EmbedObject(c).Msgf("Server never received PONG, closing connection")
			(*conn).Close()
		}

		time.Sleep(duration)
	}
}
