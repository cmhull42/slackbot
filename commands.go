package main

import (
	"encoding/json"
	"math/rand"
	"strings"
	"time"
)

// NotifyText - A non-mention message in a channel the bot has access to
func NotifyText(m Message) {
	if strings.Contains(m.Event.Text, "fuck") {
		postResponse(m.Event.Channel, "https://youtu.be/hpigjnKl7nI?t=2s")
	}
}

// NotifyMention - Called when the bot is mentioned by name
func NotifyMention(m Message) {
	if strings.Contains(m.Event.Text, "pup") {
		var j map[string]imgurresp
		puppies := imgurAPI("dog")
		json.NewDecoder(strings.NewReader(puppies)).Decode(&j)

		rand.Seed(time.Now().Unix())
		i := rand.Intn(len(j["data"].Items))

		postResponse(m.Event.Channel, j["data"].Items[i].Link)
	}
}
