package chat

import "time"

const (
	GENERAL_CHAN_PREFIX  = "#"
	LOCAL_CHAN_PREFIX    = "&"
	MODELESS_CHAN_PREFIX = "+"
)

var modes = "psimntlkbeP"

type Channel struct {
	Name            string
	Prefix          string
	Topic           string
	TopicDetails    string
	ChannelModes    string
	UserModes       []string
	Operators       map[string]*Client
	UserList        map[string]*Client
	CreateTime      time.Time
	ChannelPassword string
	Status          string
	Duration        string // whether the channel is persistent or temporary
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

	return &Channel{
		Name:            name,
		Prefix:          string(name[0]),
		Topic:           topic,
		TopicDetails:    topic_details,
		ChannelModes:    "nt", // by default only allow operators to modify topic && prevent external messages to channel (users must join it first)
		UserModes:       make([]string, 0),
		Operators:       make(map[string]*Client),
		UserList:        make(map[string]*Client),
		CreateTime:      creation_date_time,
		ChannelPassword: channel_password,
		Status:          status,
		Duration:        duration,
	}
}

func (c *Channel) SetMode(mode string) {

}

func (c *Channel) DeleteMode(mode string) {

}
