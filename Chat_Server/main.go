package main

import (
	"zirc/server"
)

func main() {

	// TODO - for whatever reason overriding go's default logger
	// 		  with "github.com/phuslu/log"'s configuration doesn't show logs
	//		  within docker containers.

	// log_level := log.DebugLevel
	// level := os.Getenv("IRC_CHAT_SERVER_LOG_LEVEL")
	// if level != "" {
	// 	log_level = log.ParseLevel(level)
	// }

	// log.DefaultLogger = log.Logger{
	// 	Level:      log_level,
	// 	Caller:     1,
	// 	TimeField:  "date",
	// 	TimeFormat: "2006-01-02",
	// 	// Writer:     &log.IOWriter{os.Stderr},
	// 	Writer: &log.ConsoleWriter{
	// 		ColorOutput:    true,
	// 		QuoteString:    true,
	// 		EndWithMessage: true,
	// 	},
	// }

	irc_server := server.NewIrcServer("", "", "", nil, nil, nil)
	irc_server.Run()
}
