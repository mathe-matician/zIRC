package chat

import "time"

const (
	PERMANENT_CHAN_PREFIX = "#"
	TEMP_CHAN_PREFIX      = "&"
	PRIVATE_CHAN_PREFIX   = "+"
)

type Channel struct {
	Name            string
	Topic           string
	TopicDetails    string
	ChannelModes    []string
	UserModes       []string
	UserList        map[string]*Client
	CreateTime      time.Time
	ChannelPassword string
	Status          string
	Duration        string // whether the channel is persistent or temporary
}

// TOPIC
// Description: A short description or subject of the channel, set by channel operators.
// Mode: Set using the +t mode to restrict who can change the topic (only operators can change it if +t is set).

// Channel Modes
//Visibility and Access:
// +p (Private): The channel is not visible in the channel list.
// +s (Secret): The channel is hidden from public view and channel lists.
// +i (Invite-Only): Users must be invited to join the channel.
// Moderation:
// +m (Moderated): Only users with voice (+v) or operator status can speak.
// +n (No External Messages): Prevents users outside the channel from sending messages to it.
// +t (Topic Protection): Only operators can change the topic.
// User Limits and Restrictions:
// +l (Limit): Sets a maximum number of users allowed in the channel.
// +k (Keyed): Requires a password (key) to join the channel.
// Bans and Exemptions:
// +b (Ban List): Bans specific users or masks from the channel.
// +e (Ban Exemption): Exempts specific users from being affected by a ban.

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
		Topic:           topic,
		TopicDetails:    topic_details,
		ChannelModes:    make([]string, 0),
		UserModes:       make([]string, 0),
		UserList:        make(map[string]*Client),
		CreateTime:      creation_date_time,
		ChannelPassword: channel_password,
		Status:          status,
		Duration:        duration,
	}
}
