package chat

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"zirc/helpers"
	rc "zirc/remote_conn"

	"github.com/google/uuid"
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

type Status struct {
	Ready bool
}

type IrcServer struct {
	SID             uuid.UUID
	DnsName         string
	Version         string
	CreationDate    time.Time
	Addr            string
	Role            string
	Description     string
	HopCount        int
	Listener        *net.Listener
	Conn            net.Conn // conn is used to keep track of s2s connections mostly
	_MessageManager *MessageManager
	_ServerManager  *ServerManager
	RoutingTable    *RoutingTable
	Servers         []*IrcServer
	Clients         []*Client
	ClientServerMap map[string]string
	ClientMap       map[string]*Client  // TODO possibly move ClientMap into the Server itself?
	ChannelMap      map[string]*Channel // TODO possibly move ChannelMap into the Server itself?
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

func NewIrcServer(cfg *Config, hop_count int, server_list *[]*IrcServer, client_list *[]*Client, config *map[string]string) *IrcServer {
	err := helpers.VerifyServerMode(cfg.Server.Server_role)
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

	enable_tls := cfg.Server.Enable_tls

	port := cfg.Server.Port
	if enable_tls {
		port = cfg.Server.Tls_port
	}

	addr := cfg.Server.Host + ":" + port

	if config == nil {
		conf := map[string]string{
			"MAX_BUFFER_SIZE":       strconv.Itoa(cfg.Server.Max_buffer_size),
			"IRC_MAX_USER_CHANNELS": strconv.Itoa(cfg.Server.Max_user_channels),
			"IRC_USER_MODES":        cfg.Server.User_modes,
			"IRC_CHANNEL_MODES":     cfg.Server.Channel_modes,
			"CAPABILITIES":          cfg.Server.Capabilities,
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

	uuid, err := uuid.NewV7()
	if err != nil {
		log.Error().Msgf("NewClient: Error generating uuid %s", err.Error())
		return nil
	}

	g_Server = &IrcServer{
		SID:             uuid,
		DnsName:         cfg.Server.Dns_name,
		Version:         cfg.Server.Server_version,
		CreationDate:    creation_date_time,
		Addr:            addr,
		Role:            cfg.Server.Server_role,
		Description:     cfg.Server.Server_description,
		HopCount:        0, // hopcount for self always == 0
		Listener:        nil,
		_MessageManager: nil,
		_ServerManager:  nil,
		RoutingTable:    NewRoutingTable(),
		Servers:         *server_list,
		Clients:         *client_list,
		ClientServerMap: make(map[string]string),
		Config:          *config,
	}

	s_manager := NewMessageManager(&g_Server.Clients, &g_Server.Servers)
	s_manager.Name = cfg.Server.Dns_name
	g_Server._MessageManager = s_manager

	server_manager := NewServerManager(nil)
	g_Server._ServerManager = server_manager

	// TODO - how do you connect to become brock_rockjaw? probably need the NickServ for this
	superadmin, session_timestamp, err := NewClient(
		"brock_rockjaw",
		"brock_rockjaw",
		nil,
		nil,
		false,
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
func GetTCPListener(enable_tls bool, tls_cert_path, tls_key_path, tls_port, irc_host, port string) net.Listener {
	var ln net.Listener
	var err error

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

// GetServerByConn finds a IrcServer in this server's server list by a net.Conn interface
func (is *IrcServer) GetServerByConn(c net.Conn) *IrcServer {
	for _, svr := range is.Servers {
		if svr.Conn == c {
			return svr
		}
	}

	return nil
}

// satisfy the Target interface to be used to identify this struct type during
// s2s communication in worker.go
func (is *IrcServer) IsTarget() {}

func (is *IrcServer) MarshalObject(e *log.Entry) {
	e.Str("dns", is.DnsName).Str("version", is.Version).Str("addr", is.Addr)
}

func (is *IrcServer) ComponentsReady() bool {
	if !is._MessageManager.ready || !is._ServerManager.ready {
		return false
	}
	return true
}

func (is *IrcServer) WaitForComponentsReady() {
	for !is.ComponentsReady() {
		log.Info().Msg("Waiting for components to be ready...")
		time.Sleep(100 * time.Millisecond)
	}
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

func (is *IrcServer) ClientMapInsert(client *Client) error {
	mm := is._MessageManager
	if mm == nil {
		err := errors.New("MessageManager is nil")
		panic(err)
	}

	err := mm.ClientMapInsert(client)
	if err != nil {
		return err
	}

	return nil
}

func (is *IrcServer) NickExists(nick string) bool {
	_, ok := is.ClientServerMap[nick]
	return ok
}

func (is *IrcServer) Run() {
	enabled_tls := G_Config.Server.Enable_tls
	tls_port := G_Config.Server.Tls_port
	addr_split := strings.Split(is.Addr, ":")
	ln := GetTCPListener(
		enabled_tls,
		G_Config.Server.Tls_cert_path,
		G_Config.Server.Tls_key_path,
		tls_port,
		addr_split[0],
		addr_split[1],
	)
	if enabled_tls {
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

		go handleConnection(conn, false)
	}
}

func (is *IrcServer) Stop() {
	// TODO - probably do some shutdown stuff before just closing all connections
	(*is.Listener).Close()
}

func handleConnection(conn net.Conn, isServer bool) {
	defer conn.Close()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	remote_addr := conn.RemoteAddr()
	remote_ip, remote_port, err := net.SplitHostPort(remote_addr.String())
	if err != nil {
		log.Error().Str("remote_addr", remote_addr.String()).Msgf("Error splitting remote addr: %s", err.Error())
	}

	var dns_name string
	// TODO
	// not always true
	if isServer {
		dns_names, err := net.LookupAddr(remote_ip)
		if err != nil {
			log.Warn().Msgf("Error performing reverse DNS lookup: %v\n", err)
		}

		log.Info().Msgf("DNS names: %v", dns_names)

		if len(dns_names) != 0 {
			// TODO
			// get all dns names?
			dns_name = dns_names[0]
		}
	}

	remote_conn := rc.NewRemoteConn(dns_name, remote_ip, remote_port)
	client, session_timestamp, err := NewClient("", "", remote_conn, conn, isServer)
	if err != nil {
		log.Error().EmbedObject(client).Msg(err.Error())
		return
	}

	// add the client to the global client list
	// log.Info().EmbedObject(client).Msgf("is.Client len before: %d", len(is.Clients))
	g_Server.Clients = append(g_Server.Clients, client)
	// log.Info().EmbedObject(client).Msgf("is.Client len after: %d", len(is.Clients))

	log.Info().EmbedObject(client).Msgf("Client connected at %s", *session_timestamp)

	var max_buffer_size int
	max_buffer_size, err = strconv.Atoi((*g_Server).Config["MAX_BUFFER_SIZE"])
	if err != nil {
		log.Error().EmbedObject(client).Msg(err.Error())
		max_buffer_size = 8192
	}

	log.Debug().EmbedObject(client).Msg("Starting forever loop for client")
	for {
		connBuffReader := bufio.NewReaderSize(conn, max_buffer_size)
		recv_buf := make([]byte, max_buffer_size)
		byteCount, err := connBuffReader.Read(recv_buf) // also ReadString('\n') but has too many edge cases

		log.Debug().EmbedObject(client).Msgf("After connBuffReader. Read %v bytes", byteCount)
		// OLD (keeping around just in case i need it)
		// block on read until the buffer has at least 1 byte.
		// just a hacky way for this to block as Read() doesn't block on its own
		// _, err := io.ReadAtLeast((*conn), recv_buf, 1)
		if err != nil {
			log.Error().EmbedObject(client).Msgf("ERR: %v", err.Error())
			if err == io.EOF {
				end_timestamp, err := client.SetSessionEndTimestamp()
				if err != nil {
					log.Error().EmbedObject(client).Msg(err.Error())
				}
				log.Info().EmbedObject(client).Msgf("Client disconnected: %s", *end_timestamp)

				// TODO
				// make these individual clean up steps a single function

				// rm client state from message manager
				g_Server._MessageManager.ClientMapDrop(client.UID.String())
				// TODO - write session duration as a metric / possible analysis

				// TODO - remove client from all channels they are on
			} else {
				log.Error().EmbedObject(client).Msgf("Error reading data from connection: %s", err.Error())
			}

			return
		}

		client_msg_start := time.Now()
		// trimmed_msg := strings.Trim(strings.TrimLeft(recv_buf, " "), "\x00")
		// OLD w/ buffer
		trimmed_msg := string(bytes.Trim(bytes.TrimLeft(recv_buf, " "), "\x00"))

		split_trimmed_msg := strings.Split(trimmed_msg, "\r\n")
		log.Debug().Msgf("split_trimmed_msg: %s", split_trimmed_msg)

		var response []byte
		for _, msg := range split_trimmed_msg {
			if len(msg) == 0 {
				continue
			}
			// TODO
			// ProcessMessage needs to pass a target not just a client (since it could be server or client)
			// ideally this is what it does...
			// we could cheap out and just have two ProcessMessage functions...
			//

			response = ProcessMessage(msg+"\r\n", client)
		}

		client_msg_end := time.Now()
		log.Info().Msgf("Client msg processing time: %v", client_msg_end.Sub(client_msg_start))

		if len(response) == 0 {
			// e.g. sometimes the server doesn't send anything back to the client
			// 		as in the case of correct password via PASS
			// reset the buffer
			continue
		}

		if strings.Contains(string(response), "ERROR") {
			log.Error().EmbedObject(client).Msgf("Critical error occurred. Closing client connection: %s", response)
			// conn closed via defer (*conn).Close() at top of func
			break
		}

		if strings.Contains(string(response), "TERMINATE") {
			log.Warn().EmbedObject(client).Msg("Terminating connection as they aren't whitelisted")
			// conn closed via defer (*conn).Close() at top of func
			break
		}

		if _, err := conn.Write(response); err != nil {
			log.Error().EmbedObject(client).Msgf("Error writing to client: %s", err.Error())
			break
		}
	}
}
