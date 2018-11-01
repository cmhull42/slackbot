package main

// Message - a slack webapi message
type Message struct {
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	EventType string `json:"type"`
	TeamID    string `json:"team_id"`
	Event     Event  `json:"event"`
}

// Event - event contents of the slack api message
type Event struct {
	Type        string `json:"type"`
	SubType     string `json:"subtype"`
	User        string `json:"user"`
	Text        string `json:"text"`
	ClientMsgID string `json:"client_msg_id"`
	Time        string `json:"ts"`
	Channel     string `json:"channel"`
	ChannelType string `json:"channel_type"`
}

// Reply - a reply through the slack api
type Reply struct {
	Text    string `json:"text"`
	Channel string `json:"channel"`
}
