package server

import (
	"io"
	"net"
	"strconv"
	"strings"

	"zirc/chat"
	c "zirc/client"
	"zirc/helpers"
	rc "zirc/remote_conn"
	sm "zirc/servermanager"

	"github.com/phuslu/log"
)

type Server interface {
	Run()
	handleConnection(conn *net.Conn)
	Stop()
}

type ServerConfig interface {
}

type IrcServer struct {
	DnsName       string
	Addr          string
	Role          string
	Listener      *net.Listener
	_ServerManger *sm.ServerManager
	Servers       []*Server
	Clients       []*c.Client
	Config        map[string]string
}

func NewIrcServer(dns_name string, addr string, server_role string, server_list *[]*Server, client_list *[]*c.Client, config *map[string]string) *IrcServer {
	if len(dns_name) == 0 {
		dns_name = helpers.GetEnv("IRC_SERVER_DNS_NAME", "localhost")
	}

	if len(server_role) == 0 {
		server_role = helpers.GetEnv("IRC_SERVER_ROLE", "leaf")
	}
	err := helpers.VerifyServerMode(server_role)
	if err != nil {
		log.Error().Msg(err.Error())
		panic(err.Error())
	}

	if server_list == nil {
		sl := make([]*Server, 0)
		server_list = &sl
	}

	if client_list == nil {
		cl := make([]*c.Client, 0)
		client_list = &cl
	}

	if len(addr) == 0 {
		addr = helpers.GetEnv("IRC_HOST", "0.0.0.0") + ":" + helpers.GetEnv("IRC_PORT", "6667")
	}

	if config == nil {
		conf := map[string]string{
			"MAX_BUFFER_SIZE": helpers.GetEnv("IRC_MAX_BUFFER_SIZE", "8192"),
		}
		config = &conf
	}

	// TODO - parse whether it should run w/ TLS or not which will determine the port used

	is := IrcServer{
		DnsName:       dns_name,
		Addr:          addr,
		Role:          server_role,
		Listener:      nil,
		_ServerManger: &sm.ServerManager{Name: dns_name},
		Servers:       *server_list,
		Clients:       *client_list,
		Config:        *config,
	}
	return &is
}

func (is *IrcServer) Run() {
	ln, err := net.Listen("tcp", is.Addr)
	if err != nil {
		log.Error().Msg(err.Error())
		panic(err.Error())
	}
	is.Listener = &ln

	log.Info().Msgf("Server started, role: %s, addr: %s, dns: %s, config: %v", is.Role, is.Addr, is.DnsName, is.Config)

	if is._ServerManger == nil {
		panic("Server's manager is null!!")
	}
	go is._ServerManger.Run()

	for {
		conn, err := (*is.Listener).Accept()
		if err != nil {
			log.Error().Msgf("Error accepting connection: %s", err.Error())
			continue
		}

		go is.handleConnection(&conn)
	}
}

func (is *IrcServer) Stop() {
	// TODO - probably do some shutdown stuff before just closing all connections
	(*is.Listener).Close()
}

func (is *IrcServer) handleConnection(conn *net.Conn) {
	defer (*conn).Close()

	remote_addr := (*conn).RemoteAddr()
	remote_ip, remote_port, err := net.SplitHostPort(remote_addr.String())
	if err != nil {
		log.Error().Str("remote_addr", remote_addr.String()).Msgf("Error splitting remote addr: %s", err.Error())
	}
	// TODO - resolve DNS name here for additional checks / verification
	// e.g. w/ servers and compare to server list
	remote_conn := rc.NewRemoteConn("", remote_ip, remote_port)
	client, session_timestamp, err := c.NewClient("", "", remote_conn)
	if err != nil {
		log.Error().EmbedObject(client).Msgf(err.Error())
		return
	}

	// add the client to the global client list
	is.Clients = append(is.Clients, client)

	log.Info().EmbedObject(client).Msgf("Client connected at %s", *session_timestamp)

	var max_buffer_size int
	max_buffer_size, err = strconv.Atoi((*is).Config["MAX_BUFFER_SIZE"])
	if err != nil {
		log.Error().EmbedObject(client).Msgf(err.Error())
		max_buffer_size = 8192
	}
	recv_buf := make([]byte, max_buffer_size)

	for {
		// TODO - clear buffers so no extra data is sent?
		// TODO - send periodic PING commands
		//		  if no PONG is received, terminate the connection
		//		  used to determine dead connections
		_, err := (*conn).Read(recv_buf)
		if err != nil {
			if err == io.EOF {
				end_timestamp, err := client.SetSessionEndTimestamp()
				if err != nil {
					log.Error().EmbedObject(client).Msg(err.Error())
				}
				log.Info().EmbedObject(client).Msgf("Client disconnected: %s", *end_timestamp)
				// TODO - write session duration as a metric / possible analysis
			} else {
				log.Error().EmbedObject(client).Msgf("Error reading data from connection: %s", err.Error())
			}
			return
		}

		response := chat.ProcessMessage(&recv_buf, client, is._ServerManger)

		if len(response) == 0 {
			// e.g. sometimes the server doesn't send anything back to the client
			// 		as in the case of correct password via PASS
			continue
		}

		if _, err := (*conn).Write(response); err != nil {
			log.Error().EmbedObject(client).Msgf("Error writing to client: %s", err.Error())
			break
		}

		if strings.Contains(string(response), "ERROR") {
			log.Error().EmbedObject(client).Msgf("Critical error occurred. Closing client connection: %s", response)
			(*conn).Close()
			break
		}
	}
}
