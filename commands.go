package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

var eightballResponses = [...]string{"It is certain.", "As I see it, yes", "Reply hazy, try again.", "Don't count on it.", "It is decidedly so.", "Most likely", "Ask again later.", "My reply is no.", "Without a doubt.", "Outlook good.", "Better not tell you now.", "My sources say no.", "Yes - definitely.", "Yes.", "Cannot predict now.", "Outlook not so good.", "You may rely on it.", "Signs point to yes.", "Concentrate and ask again.", "Very doubtful."}

type commands struct {
	Messager messager
}

// NotifyText - A non-mention message in a channel the bot has access to
func (c commands) NotifyText(m Message) {
	rand.Seed(time.Now().Unix())
	i := rand.Intn(500)
	log.Printf("Rolling the dice and got %d", i)
	if i == 69 {
		c.Messager.postResponse(m.Event.Channel, "That's a microaggression. Reported.")
	}
}

// NotifyMention - Called when the bot is mentioned by name
func (c commands) NotifyMention(m Message) {

	imgurRegex := regexp.MustCompile("pic(?:ture)? of (?:a )? ?([a-zA-Z ]+)")
	if imgurRegex.MatchString(m.Event.Text) {
		thing := imgurRegex.FindStringSubmatch(m.Event.Text)[1]

		log.Printf("]" + thing + "[")
		if thing == "aw heck" {
			c.Messager.postResponse(m.Event.Channel, "Please don't make me look that up again")
			return
		}

		var j map[string][]imgurresp
		things := imgurAPI(thing)
		json.NewDecoder(strings.NewReader(things)).Decode(&j)

		rand.Seed(time.Now().Unix())

		if len(j["data"]) == 0 {
			c.Messager.postResponse(m.Event.Channel, "No results found.")
		} else {
			i := rand.Intn(len(j["data"]))

			if len(j["data"][i].Images) == 0 {
				c.Messager.postResponse(m.Event.Channel, j["data"][i].Link)
			} else {
				c.Messager.postResponse(m.Event.Channel, j["data"][i].Images[0].Link)
			}
		}
	}

	issueRegex := regexp.MustCompile(`issue (\d+) ?`)
	if issueRegex.MatchString(m.Event.Text) {
		num := issueRegex.FindStringSubmatch(m.Event.Text)[1]
		c.Messager.postResponse(m.Event.Channel, "http://dev-tracker/Lists/All%20Suite%20Issues/DispForm.aspx?ID="+num)
	}

	eightballRegex := regexp.MustCompile(`.+\?$`)
	if eightballRegex.MatchString(m.Event.Text) {
		rand.Seed(time.Now().Unix())
		i := rand.Intn(20)
		c.Messager.postResponse(m.Event.Channel, eightballResponses[i])
	}
}
