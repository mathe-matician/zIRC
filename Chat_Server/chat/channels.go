package chat

import (
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/phuslu/log"
)

const (
	GENERAL_CHAN_PREFIX  = "#"
	LOCAL_CHAN_PREFIX    = "&"
	MODELESS_CHAN_PREFIX = "+"
	MAX_CHAN_LEN         = 100
)

var modes = "psimntlkbeP"

type Mode struct {
	ModeChar string
	Params   string
}

type Channel struct {
	Name            string
	Prefix          string
	Topic           string
	TopicDetails    string
	ChannelModes    []Mode
	UserModes       []string
	Operators       map[string]*Client
	UserList        map[string]*Client
	CreateTime      time.Time
	ChannelPassword string
	Status          string
	Duration        string // whether the channel is persistent or temporary
	mu              sync.Mutex
}

// TOPIC
// Description: A short description or subject of the channel, set by channel operators.
// Mode: Set using the +t mode to restrict who can change the topic (only operators can change it if +t is set).

// USER MODES
// +o (Operator): Grants the user operator status, allowing them to manage the channel.
// +v (Voice): Allows the user to speak in a moderated channel.

// CHANNEL STATUS
// Active/Inactive: Reflects whether the channel is currently active (has users) or inactive (empty and possibly removed if temporary).

func NewChannel(name, topic, topic_details, channel_password, status, duration string) *Channel {
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

	default_channel_modes := []Mode{
		{ModeChar: "n", Params: ""},
		{ModeChar: "t", Params: ""},
	}

	return &Channel{
		Name:         name,
		Prefix:       string(name[0]),
		Topic:        topic,
		TopicDetails: topic_details,
		// by default only allow operators to modify topic && prevent external messages to channel (users must join it first)
		ChannelModes:    default_channel_modes,
		UserModes:       make([]string, 0),
		Operators:       make(map[string]*Client),
		UserList:        make(map[string]*Client),
		CreateTime:      creation_date_time,
		ChannelPassword: channel_password,
		Status:          status,
		Duration:        duration,
	}
}

func (c *Channel) IsTarget() {}

func (c *Channel) FmtModes() string {
	modes := ""
	params := ""
	for _, m := range c.ChannelModes {
		modes += m.ModeChar
		params += " " + m.Params
	}
	return modes + " " + strings.TrimLeft(params, " ")
}

func NewMode(mode string, params string) *Mode {
	return &Mode{
		ModeChar: mode,
		Params:   params,
	}
}

func (c *Channel) AddMode(mode Mode) {
	c.mu.Lock()
	defer c.mu.Unlock()
	char := mode.ModeChar
	if (char == "k" && c.HasMode(char)) || (!mode_requires_params(char) && c.HasMode(char)) {
		// even though k takes params, only have 1 k in list (i.e. having multiple passwords doesn't make sense)
		return
	}

	c.ChannelModes = append(c.ChannelModes, mode)
}

func (c *Channel) HasModeAndValue(mode, params string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, m := range c.ChannelModes {
		if m.ModeChar == mode && m.Params == params {
			return true
		}
	}

	return false
}

func (c *Channel) HasMode(mode string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, m := range c.ChannelModes {
		if m.ModeChar == mode {
			return true
		}
	}
	return false
}

func (c *Channel) RemoveMode(_mode string, param string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// if !strings.Contains(c.ChannelModes, mode) {
	// 	return
	// }
	// c.ChannelModes = strings.Replace(c.ChannelModes, mode, "", 1)
	m := NewMode(_mode, param)
	mode := *m
	for i, m := range c.ChannelModes {
		// l and k modes don't need params to remove them
		// so don't worry about those params if we are removing them as they will be empty
		if ((m.ModeChar == "l" || m.ModeChar == "k") && m.ModeChar == mode.ModeChar) || (m.ModeChar == mode.ModeChar && m.Params == mode.Params) {
			if i+1 > len(c.ChannelModes) {
				log.Error().Msgf("Error when removing mode %v from channelmodes: %v", mode, c.ChannelModes)
				break
			}
			c.ChannelModes = slices.Delete(c.ChannelModes, i, i+1)
			return
		}
	}

	log.Debug().Msgf("Mode %v wasn't found in channelmodes: %v", mode, c.ChannelModes)
}

// chans can't have ' ', escape characters, prefix characters,
// chans must be A-Z, a-z, 0-9, and can have any of these: - _ . /
var chan_re = regexp.MustCompile(`^[#&+][A-Za-z0-9_\-\.\/]*$`)

func IsValidChanName(channel string) bool {
	if len(chan_re.FindAllStringSubmatch(channel, -1)) == 0 || len(channel) > MAX_CHAN_LEN {
		return false
	}
	return true
}
