package chat

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
	"zirc/helpers"

	"github.com/phuslu/log"
	"gopkg.in/yaml.v3"
)

var server_manager_commands = map[string]Command{
	"PASS":     *NewCommand(pass, map[string]string{"cap_req": "sasl"}, true),
	"NETINFO":  *NewCommand(netinfo, map[string]string{"cap_req": "sasl"}, true),
	"UID":      *NewCommand(not_implemented, map[string]string{"cap_req": "sasl"}, true),
	"CAPAB":    *NewCommand(not_implemented, map[string]string{"cap_req": "sasl"}, true),
	"SERVER":   *NewCommand(server, map[string]string{"cap_req": "sasl"}, true),
	"CAP":      *NewCommand(not_implemented, map[string]string{"cap_req": "sasl"}, true),
	"SJOIN":    *NewCommand(not_implemented, map[string]string{"cap_req": "sasl"}, true),
	"PING":     *NewCommand(not_implemented, map[string]string{"cap_req": "sasl"}, true),
	"PONG":     *NewCommand(not_implemented, map[string]string{"cap_req": "sasl"}, true),
	"SQUIT":    *NewCommand(not_implemented, map[string]string{"cap_req": "sasl"}, true),
	"CONNECT":  *NewCommand(not_implemented, map[string]string{"cap_req": "sasl"}, true),
	"BURST":    *NewCommand(not_implemented, map[string]string{"cap_req": "sasl"}, true),
	"ENDBURST": *NewCommand(not_implemented, map[string]string{"cap_req": "sasl"}, true),
	"EUID":     *NewCommand(not_implemented, map[string]string{"cap_req": "sasl"}, true),
	//	b. BURST / EUID (IRCv3)
	//
	// Purpose: Synchronize state after a netsplit or during initial connection.
	// Routing Table Use:
	// The server ensures that its routing table is updated with the correct paths for users and channels.
}

func serverCommandValidation(cmd string) (*Command, Response) {
	// TODO
	// should we check for s2s password again?
	// probably not since we could only get here if we had a valid s2s connection
	// i.e. we already authenticated in Connect
	log.Debug().Msgf("Server CMD: %s", cmd)
	return_cmd, ok := server_manager_commands[cmd]
	if !ok {
		return nil, ERR_UNKNOWNCOMMAND("")
	}

	res := Reply{
		code: "-1",
		msg:  "hi",
	}
	return &return_cmd, &res
}

type ServerConnection struct {
	Name         string `yaml:"name"` // IRC server name used as index in routing table
	Host         string `yaml:"host"` // generally used as DNS name for preconfigured servers
	Ip           string `yaml:"ip"`   // the ip address of the server
	Port         string `yaml:"port"`
	AutoConnect  bool   `yaml:"auto_connect"`  // for preconfigured servers, if true auto connect to this server on startup
	PasswordFile string `yaml:"password_file"` // the path on the server to read a password file from
	Password     string
	// Conn         net.Conn
}

type ServerManagerConfig struct {
	ServerList *map[string]ServerConnection
}

type ServerManager struct {
	Addr              string
	serverIpWhitelist map[string]string // TODO, what should the value be here?
	Config            *ServerManagerConfig
	ready             bool
}

func NewServerManager() *ServerManager {
	sm_config := NewServerManagerConfig()
	log.Debug().Msgf("Serverlist: %v", sm_config.ServerList)
	sm := ServerManager{
		Addr:              G_Config.S2S.Port,
		Config:            sm_config,
		serverIpWhitelist: make(map[string]string),
	}

	if len(*sm.Config.ServerList) != 0 {
		for _, server := range *sm.Config.ServerList {
			// prefer mapping by IP, but that may not always be the case
			// there may be better ways to do this.
			if server.Ip != "" {
				// tmp because idk why the value is a string
				sm.serverIpWhitelist[server.Ip] = "whitelisted"
			} else if server.Host != "" {
				sm.serverIpWhitelist[server.Host] = "whitelisted"
			}

			if server.AutoConnect {
				// TODO
				// need to account for race condition here
				// g_Server.WaitForComponentsReady()
				go Connect(server)
			}
		}

		// wait for any failures or successful connections?
		// select {
		// case pong_token := <-c.PingPongChan:
		// 	pongToken := pong_token[1:]
		// 	if pongToken != pingToken {
		// 		log.Info().EmbedObject(c).Msgf("%s != %s", pongToken, pingToken)
		// 		log.Info().EmbedObject(c).Msgf("%s len: %d", pongToken, len(pongToken))
		// 		log.Info().EmbedObject(c).Msgf("%s len: %d", pingToken, len(pingToken))
		// 		(*conn).Close()
		// 	}
		// 	log.Info().EmbedObject(c).Msgf("PONG successful")
		// case <-time.After(time.Duration(timeout) * time.Second):
		// 	log.Info().EmbedObject(c).Msgf("Server never received PONG, closing connection")
		// 	(*conn).Close()
		// }
	}

	return &sm
}

type ServerConnectionPkg struct {
	ServerName       string
	ConnectedServers []ServerConnection // todo this probably shouldn't contain other server passwords
	ClientList       []Client
	ChannelList      []Channel
}

// connect to another server
func Connect(connection ServerConnection) {
	// - server name
	// - directly connected servers
	// - user and channel lists it manages
	// - network paths to reach other servers
	// e.g.
	//       A -- B -- C
	//			  |
	//			  D
	// When B connects to A, it says "I'm connced to C and D"
	// A could then hold:
	// []Servers{
	//		[]B {C, D}
	//  }
	// If multiple paths exist, then path selection considers:
	// 		- shortest path (least hops)
	//		- link cost (if the network assigns weights to connections)

	//
	// routing table updates
	//

	// Final Command Sequence (Summary)
	// PASS <password> <protocol_version> <flags>
	// SERVER <servername> <hopcount> :<description>
	// NETINFO (If used)
	// BURST (If TS6 is used)
	// NICK (For every active user)
	// JOIN (For every active channel)
	// MODE (Sync channel modes)
	// TOPIC (Sync channel topics)
	// SJOIN (If using TS6)
	// ENDBURST (If required)
	connect_start_time := time.Now()
	log.Debug().Msg("Connect start")

	addr := connection.Host + ":" + connection.Port
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Error().Msgf("Failed to connect to server: %s", addr)
		return
	}
	defer conn.Close()

	CRLF := "\r\n"
	passFlags := " "
	pass := fmt.Sprintf("PASS %s %s%s%s", connection.Password, G_Config.Server.Irc_verison, passFlags, CRLF)

	// TODO
	// abstract this into its own command
	// see server_cmd.go
	hopcount := 1 // hopcount 1 since Connect() will always be a direct connection to another server
	server := fmt.Sprintf("SERVER %s %d :A test server %s", G_Config.Server.Server_name, hopcount, CRLF)
	// netinfo := ":server1.example.com NETINFO 1707500000 1707500001 0 J10 TS6 6 :server1.example.com"
	// Breakdown:
	// 	:server1.example.com → The source server sending the NETINFO.
	// 	NETINFO → The command name.
	// 	1707500000 → The network's creation timestamp (typically a Unix timestamp).
	// 	1707500001 → The current timestamp when the message is sent.
	// 	0 → The protocol version (0 for TS6).
	// 	J10 → The network's numeric version identifier (implementation-specific).
	// 	TS6 → The timestamping protocol in use (TS6 in this case).
	// 	6 → The number of additional parameters.
	// 	:server1.example.com → The server name that originated the message.

	// TODO
	// how to handle admin user brockrockjaw?
	// they should be excluded as they will exist across all servers
	// SYNC all clients on this server
	// if len(*g_Server._MessageManager.ClientList) != 0 {
	// 	for _, c := range *g_Server._MessageManager.ClientList {
	// 		msg := fmt.Sprintf("NICK %s %v %s %s %s %s :%s %s", c.Nick(), c.NickTimestamp, g_Server.DnsName, c.User(), c.Host, g_Server.DnsName, c.RealName, CRLF)
	// 		_, err = conn.Write([]byte(msg))
	// 		if err != nil {
	// 			log.Error().Msgf("Error writing nicks to server: %s", err.Error())
	// 			return
	// 		}
	// 	}
	// }

	msg := fmt.Sprintf("%s%s", pass, server)
	// burst := ""
	// sync commands
	// nick
	//

	_, err = conn.Write([]byte(msg))
	if err != nil {
		log.Error().Msgf("Error writing: %s", err.Error())
		return
	}

	connect_end_time := time.Now()
	log.Debug().Msgf("Server handshake successful. Took: %v", connect_start_time.Sub(connect_end_time))

	// need to whitelist the IP upon connection for our server connecting to the other
	// wl := *whiteList
	// ip, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	// if err != nil {
	// 	log.Error().Msgf("error when connecting to remote server: %s", err.Error())
	// 	return
	// }
	// wl[ip] = connection
}

func NewServerConnection() *ServerConnection {
	return &ServerConnection{}
}

func NewServerManagerConfig() *ServerManagerConfig {
	server_list := make([]ServerConnection, 0)
	server_map := make(map[string]ServerConnection)
	config_path := helpers.GetEnv("IRC_S2S_CONFIG_FILE", "/chat_server/s2s_config.yaml")
	config_data, err := os.ReadFile(config_path)
	if err != nil {
		// probaby doens't have to be an error since we don't _have_ to have a config here
		// as of right now it only loads connected servers
		// _could_ support more in the future, though
		log.Error().Msgf("Error reading config file %v", err)
		return &ServerManagerConfig{
			ServerList: &server_map,
		}
	}

	// TODO
	// check for zero length config_data
	// no need to unmarshal then as no additional servers exist
	if len(config_data) == 0 {
		log.Info().Msgf("No data in server manager config")
		return &ServerManagerConfig{
			ServerList: &server_map,
		}
	}

	err = yaml.Unmarshal([]byte(config_data), &server_list)
	if err != nil {
		log.Error().Msgf("ServerManagerConfig couldnt unmarshal config: %v", err)
		panic(err)
	}

	for _, server := range server_list {
		if server.Ip != "" {
			server_map[server.Ip] = server
		}
		server_map[server.Ip] = server
		if server.PasswordFile == "" {
			continue
		}

		server_password, err := os.ReadFile(server.PasswordFile)
		if err != nil {
			log.Error().Msgf("Error reading %s password file: %v", server.Host, err)
			continue
		}

		server.Password = string(server_password)
	}

	return &ServerManagerConfig{
		ServerList: &server_map,
	}
}

func (sm *ServerManager) init() {

}

func (sm *ServerManager) AddToServerIpWhiteList(ip, args string) {
	sm.serverIpWhitelist[ip] = args
}

func (sm *ServerManager) RemoveFromServerIpWhiteList(ip string) {
	delete(sm.serverIpWhitelist, ip)
}

// Allows this server to Join an existing IRC network
// TODO
// should this be an init function?
// e.g. get Join config
func (sm *ServerManager) Join(ip string) {
	/*
		1. PASS
		2. CAPAB (capability announcement)
			The server initiating the connection sends a CAPAB (capability) message to the receiving server. This message lists the features or protocol extensions that the server supports.

			CAPAB :QS EX IE KLN UNKLN

			Here, the server is advertising that it supports several features:

			QS: Quiet Channel Synchronization
			EX: Extended bans
			IE: Invite exceptions
			KLN: Kill line commands
			UNKLN: Undo kill line commands

			Each capability represents a specific feature or extension that allows enhanced communication between the servers.

		3. CAPAB matching
			Once the first server announces its capabilities, the receiving server also sends back a CAPAB message indicating its own supported capabilities.

			The two servers then compare their lists of supported capabilities. Both servers will use the features they both support for further communication. If a capability is not supported by one of the servers, it is not used in future communication.

		4. Negotiate features

			If both servers support advanced features (like extended bans, invite exceptions, etc.), they will use these capabilities in communication. For instance, if both servers support EX (extended bans), they will use extended ban syntax when communicating about bans across the network.

			If one of the servers does not support a specific feature, the two servers will revert to the basic, standard IRC protocol features for that area of communication.

		5. SERVER cmd

		https://docs.inspircd.org/server/examples/connection/
	*/

}

// Checks to see whether the server trying to connect to this network
// has entered in a valid password via the PASS command
func (sm *ServerManager) ValidS2SPassword(password string) bool {
	// check db first
	// if can't connect to db check locally for conf file or env variable
	// if all fails, then don't let join

	// last fallback is env var
	return G_Config.S2S.Password_file == password
}

func (sm *ServerManager) Run() {
	// listens on port 7000 for servers joining the irc network
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()
	ln := GetTCPListener(
		G_Config.S2S.Enable_tls,
		G_Config.S2S.Tls_cert_path,
		G_Config.S2S.Tls_key_path,
		G_Config.S2S.Tls_port,
		G_Config.Server.Host,
		G_Config.S2S.Port,
	)

	log.Info().Msgf("ServerManager listening: %s", sm.Addr)
	sm.ready = true

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Error().Msgf("Error accepting connection: %s", err.Error())
			continue
		}

		if conn.RemoteAddr().Network() != "tcp" {
			conn.Close()
			continue
		}

		if !sm.whiteListedServerIp(conn.RemoteAddr().String()) {
			conn.Close()
			continue
		}

		go handleConnection(&conn, true)

		// go sm.handleConnection(&conn)
	}
}

// Checks to see if the connected IP is an expected server IP
// i.e. do we expect this IP to be communicating with the ServerManager
// Note: arg ip is expected to be ip:port actually.
func (sm *ServerManager) whiteListedServerIp(ip string) bool {
	server_ip := ip
	var err error
	if strings.Contains(server_ip, ":") {
		server_ip, _, err = net.SplitHostPort(ip)
		if err != nil {
			log.Error().Msgf("ServerManager Error extracting IP: %s", err)
			return false
		}
	}

	if _, ok := sm.serverIpWhitelist[server_ip]; !ok {
		log.Error().Msgf("ServerManager not a whitelisted server ip!: %s", server_ip)
		return false
	}
	return true
}

func (sm *ServerManager) handleConnection(conn *net.Conn) {
	defer (*conn).Close()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ServerManager recovered from panic:", r)
		}
	}()
	remote_addr := (*conn).RemoteAddr()
	log.Info().Msgf("Handling ServerManager connection: %s", remote_addr)

	// TODO
	// does this size differ for S2S?
	max_buffer_size := G_Config.S2S.Max_buffer_size

	for {
		if !sm.whiteListedServerIp((*conn).RemoteAddr().String()) {
			// need to always check if this IP is still valid as we handle this connection
			// as we could update the whitelist mid connection
			break
		}

		// block on read until the buffer has at least 1 byte.
		// just a hacky way for this to block as Read() doesn't block on its own
		recv_buf := make([]byte, max_buffer_size)
		_, err := io.ReadAtLeast((*conn), recv_buf, 1)
		if err != nil {
			if err == io.EOF {
				log.Info().Msgf("Server %s disconnected", remote_addr)
			}
			return
		}

		if _, err := (*conn).Write([]byte("")); err != nil {
			break
		}
	}
}

func (sm *ServerManager) Route(msg []byte, server string, hopcount int) {
	log.Debug().Msgf("ServerManager::Route::msg: %s", msg)
	// TODO
	// we already have a connection to all direct servers
	// we need to now send these details via some channel
	// since those connections are running in their own go routines.
}
