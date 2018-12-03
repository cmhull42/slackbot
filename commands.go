package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

// NotifyText - A non-mention message in a channel the bot has access to
func NotifyText(m Message) {
	rand.Seed(time.Now().Unix())
	i := rand.Intn(500)
	log.Printf("Rolling the dice and got %d", i)
	if i == 69 {
		postResponse(m.Event.Channel, "That's a microaggression. Reported.")
	}
}

// NotifyMention - Called when the bot is mentioned by name
func NotifyMention(m Message) {
	if strings.Contains(m.Event.Text, "pup") || strings.Contains(m.Event.Text, "dog") {
		var j map[string]imgurresp
		puppies := imgurAPI("dog")
		json.NewDecoder(strings.NewReader(puppies)).Decode(&j)

		rand.Seed(time.Now().Unix())
		i := rand.Intn(len(j["data"].Items))

		postResponse(m.Event.Channel, j["data"].Items[i].Link)
	}

	issueRegex := regexp.MustCompile(`issue (\d+) ?`)
	if issueRegex.MatchString(m.Event.Text) {
		num := issueRegex.FindStringSubmatch(m.Event.Text)[1]
		postResponse(m.Event.Channel, "http://dev-tracker/Lists/All%20Suite%20Issues/DispForm.aspx?ID="+num)
	}
}
