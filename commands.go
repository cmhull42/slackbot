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

	imgurRegex := regexp.MustCompile("pic(?:ture)? of (?:a)? ?([a-zA-Z ]+)")
	if imgurRegex.MatchString(m.Event.Text) {
		thing := imgurRegex.FindStringSubmatch(m.Event.Text)[1]
		var j map[string][]imgurresp
		things := imgurAPI(thing)
		json.NewDecoder(strings.NewReader(things)).Decode(&j)

		rand.Seed(time.Now().Unix())

		if len(j["data"]) == 0 {
			postResponse(m.Event.Channel, "No results found.")
		} else {
			i := rand.Intn(len(j["data"]))

			if len(j["data"][i].Images) == 0 {
				postResponse(m.Event.Channel, j["data"][i].Link)
			} else {
				postResponse(m.Event.Channel, j["data"][i].Images[0].Link)
			}
		}
	}

	issueRegex := regexp.MustCompile(`issue (\d+) ?`)
	if issueRegex.MatchString(m.Event.Text) {
		num := issueRegex.FindStringSubmatch(m.Event.Text)[1]
		postResponse(m.Event.Channel, "http://dev-tracker/Lists/All%20Suite%20Issues/DispForm.aspx?ID="+num)
	}
}
