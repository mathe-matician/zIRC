package chat

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"zirc/helpers"
	rc "zirc/remote_conn"

	"github.com/phuslu/log"
)

var g_Server *IrcServer

type Server interface {
	Run()
	handleConnection(conn *net.Conn)
	Stop()
}

type ServerConfig interface {
}

type IrcServer struct {
	DnsName         string
	Version         string
	CreationDate    time.Time
	Addr            string
	Role            string
	Listener        *net.Listener
	_MessageManager *MessageManager
	_ServerManager  *ServerManager
	Servers         []*IrcServer
	Clients         []*Client
	Config          map[string]string
}

// type Task struct {
// 	id     uuid.UUID
// 	Type   string
// 	weight float64
// 	task   string
// }

// task_weight_mapping := map[string]float64 {

// }

func NewIrcServer(dns_name string, version string, addr string, server_role string, server_list *[]*IrcServer, client_list *[]*Client, config *map[string]string) *IrcServer {
	if len(dns_name) == 0 {
		dns_name = helpers.GetEnv("IRC_SERVER_DNS_NAME", "localhost")
	}

	if len(version) == 0 {
		version = helpers.GetEnv("IRC_SERVER_VERSION", "v99.99.99+default")
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
		sl := make([]*IrcServer, 0)
		server_list = &sl
	}

	if client_list == nil {
		cl := make([]*Client, 0)
		client_list = &cl
	}

	enable_tls, err := strconv.ParseBool(helpers.GetEnv("IRC_ENABLE_TLS", "false"))
	if err != nil {
		log.Error().Msg(err.Error())
	}
	port := helpers.GetEnv("IRC_PORT", "6667")
	if enable_tls {
		port = helpers.GetEnv("IRC_TLS_PORT", "6697")
	}

	if len(addr) == 0 {
		addr = helpers.GetEnv("IRC_HOST", "0.0.0.0") + ":" + port
	}

	if config == nil {
		conf := map[string]string{
			"MAX_BUFFER_SIZE":       helpers.GetEnv("IRC_MAX_BUFFER_SIZE", "8192"),
			"IRC_MAX_USER_CHANNELS": helpers.GetEnv("IRC_MAX_USER_CHANNELS", "20"),
			"IRC_USER_MODES":        helpers.GetEnv("IRC_USER_MODES", "oiws"),
			"IRC_CHANNEL_MODES":     helpers.GetEnv("IRC_CHANNEL_MODES", "opsmt"),
			"CAPABILITIES":          helpers.GetEnv("IRC_SERVER_CAPABILITIES", ""),
		}
		config = &conf
	}

	// TODO - parse whether it should run w/ TLS or not which will determine the port used
	now := time.Now()

	creation_date_time := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
		now.Nanosecond(),
		now.Location(),
	)

	g_Server = &IrcServer{
		DnsName:         dns_name,
		Version:         version,
		CreationDate:    creation_date_time,
		Addr:            addr,
		Role:            server_role,
		Listener:        nil,
		_MessageManager: nil,
		_ServerManager:  nil,
		Servers:         *server_list,
		Clients:         *client_list,
		Config:          *config,
	}

	s_manager := NewMessageManager(&g_Server.Clients, &g_Server.Servers)
	s_manager.Name = dns_name
	g_Server._MessageManager = s_manager

	server_manager := NewServerManager()
	g_Server._ServerManager = server_manager

	// TODO - how do you connect to become brock_rockjaw? probably need the NickServ for this
	superadmin, session_timestamp, err := NewClient(
		"brock_rockjaw",
		"brock_rockjaw",
		nil,
		nil,
	)

	if err != nil {
		panic("cant create defaut admin user")
	}

	*g_Server._MessageManager.ClientList = append(*g_Server._MessageManager.ClientList, superadmin)
	log.Info().Msgf("Default superadmin brock_rockjaw created at %s", *session_timestamp)

	return g_Server
}

// GetTCPListener checks to see if TLS is enabled for the particular listener
// if it is, it creates a TLS Listener with the provided TLS env vars
// else it returns a normal Listener
func GetTCPListener(enableTls, tls_cert_path, tls_key_path, tls_port, irc_host, port string) net.Listener {
	var ln net.Listener
	var err error
	enable_tls, err1 := strconv.ParseBool(enableTls)
	if err1 != nil {
		log.Error().Msg(err1.Error())
	}

	if enable_tls {
		config := helpers.CreateTLSConfig(tls_cert_path, tls_key_path)
		tls_port := tls_port
		ln, err = tls.Listen("tcp", irc_host+":"+tls_port, config)
		if err != nil {
			panic(err)
		}
	} else {
		ln, err = net.Listen("tcp", irc_host+":"+port)
		if err != nil {
			log.Error().Msg(err.Error())
			panic(err.Error())
		}
	}
	return ln
}

func (is *IrcServer) HasCapability(cap string) bool {
	caps := strings.Split(is.Config["CAPABILITIES"], " ")
	for _, c := range caps {
		if c == cap {
			return true
		}
	}
	return false
}

func (is *IrcServer) GetChannelMap() *map[string]*Channel {
	return is._MessageManager.ChannelMap
}

func (is *IrcServer) GetClientMap() *map[string]*Client {
	return is._MessageManager.ClientMap
}

func (is *IrcServer) Run() {
	enabled_tls := helpers.GetEnv("IRC_ENABLE_TLS", "false")
	tls_port := helpers.GetEnv("IRC_TLS_PORT", "6697")
	addr_split := strings.Split(is.Addr, ":")
	ln := GetTCPListener(
		enabled_tls,
		helpers.GetEnv("IRC_TLS_CERT_PATH", "./.tls/server.crt"),
		helpers.GetEnv("IRC_TLS_KEY_PATH", "./.tls/server.key"),
		tls_port,
		addr_split[0],
		addr_split[1],
	)
	if enabled_tls == "true" {
		is.Addr = tls_port
	}
	is.Listener = &ln

	log.Info().Msgf("Server started, role: %s, addr: %s, dns: %s, config: %v", is.Role, is.Addr, is.DnsName, is.Config)

	if is._MessageManager == nil {
		panic("Server's manager is null!!")
	}
	go is._MessageManager.Run()
	go is._ServerManager.Run() // handler for servers trying to join the network

	for {
		conn, err := (*is.Listener).Accept()
		if err != nil {
			log.Error().Msgf("Error accepting connection: %s", err.Error())
			continue
		}

		if conn.RemoteAddr().Network() != "tcp" {
			conn.Close()
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
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	remote_addr := (*conn).RemoteAddr()
	remote_ip, remote_port, err := net.SplitHostPort(remote_addr.String())
	if err != nil {
		log.Error().Str("remote_addr", remote_addr.String()).Msgf("Error splitting remote addr: %s", err.Error())
	}
	// TODO - resolve DNS name here for additional checks / verification
	// e.g. w/ servers and compare to server list
	remote_conn := rc.NewRemoteConn("", remote_ip, remote_port)
	client, session_timestamp, err := NewClient("", "", remote_conn, conn)
	if err != nil {
		log.Error().EmbedObject(client).Msgf(err.Error())
		return
	}

	// add the client to the global client list
	// log.Info().EmbedObject(client).Msgf("is.Client len before: %d", len(is.Clients))
	is.Clients = append(is.Clients, client)
	// log.Info().EmbedObject(client).Msgf("is.Client len after: %d", len(is.Clients))

	log.Info().EmbedObject(client).Msgf("Client connected at %s", *session_timestamp)

	var max_buffer_size int
	max_buffer_size, err = strconv.Atoi((*is).Config["MAX_BUFFER_SIZE"])
	if err != nil {
		log.Error().EmbedObject(client).Msgf(err.Error())
		max_buffer_size = 8192
	}

	for {
		// TODO - clear buffers so no extra data is sent?
		// TODO - send periodic PING commands
		//		  if no PONG is received, terminate the connection
		//		  used to determine dead connections

		// block on read until the buffer has at least 1 byte.
		// just a hacky way for this to block as Read() doesn't block on its own
		recv_buf := make([]byte, max_buffer_size)
		_, err := io.ReadAtLeast((*conn), recv_buf, 1)
		// _, err := (*conn).Read(recv_buf)
		if err != nil {
			if err == io.EOF {
				end_timestamp, err := client.SetSessionEndTimestamp()
				if err != nil {
					log.Error().EmbedObject(client).Msg(err.Error())
				}
				log.Info().EmbedObject(client).Msgf("Client disconnected: %s", *end_timestamp)
				// TODO - remove client state from MessageManager!!!
				// TODO - write session duration as a metric / possible analysis

				// TODO - remove client from all channels they are on
			} else {
				log.Error().EmbedObject(client).Msgf("Error reading data from connection: %s", err.Error())
			}

			//if err == io.ErrShortBuffer
			return
		}

		server_metadata := map[string]string{
			"name":         is.DnsName,
			"version":      is.Version,
			"date":         is.CreationDate.String(),
			"usermodes":    is.Config["IRC_USER_MODES"],
			"channelmodes": is.Config["IRC_CHANNEL_MODES"],
		}

		trimmed_msg := string(bytes.Trim(bytes.TrimLeft(recv_buf, " "), "\x00"))

		split_trimmed_msg := strings.Split(trimmed_msg, "\r\n")
		log.Debug().Msgf("split_trimmed_msg: %s", split_trimmed_msg)

		var response []byte
		for _, msg := range split_trimmed_msg {
			if len(msg) == 0 {
				continue
			}
			response = ProcessMessage(msg+"\r\n", client, is._MessageManager.Task_runner, server_metadata, is._MessageManager.ChannelMap)
		}

		if len(response) == 0 {
			// e.g. sometimes the server doesn't send anything back to the client
			// 		as in the case of correct password via PASS
			// reset the buffer
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
