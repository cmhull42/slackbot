package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var config configuration
var c commands

func main() {

	file, err := os.Open("conf.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	var m messager
	m = slackMessager{Config: config}

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		if args[i] == "--test" {
			m = testMessager{Config: config}
		}
	}

	c = commands{Messager: m}

	router := mux.NewRouter()
	router.HandleFunc("/", postMessage).Methods("POST")
	router.HandleFunc("/say", receiveMessage).Methods("POST")
	log.Fatal(http.ListenAndServe(":9803", router))
}

func postMessage(w http.ResponseWriter, r *http.Request) {
	var m Message
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &m)

	log.Print(string(b))
	switch m.EventType {
	case "url_verification":
		urlVerification(w, r, m)
	case "event_callback":
		go eventCallback(w, r, m)
	}
}

func receiveMessage(w http.ResponseWriter, r *http.Request) {
	var s sayMessage
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &s)

	log.Print(string(b))
	if s.Challenge == config.SayChallenge {
		c.Messager.postResponse(s.Channel, s.Message)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func urlVerification(w http.ResponseWriter, r *http.Request, m Message) {
	if m.Token == config.VerificationToken {
		challengeResponse := challengeResponse{Challenge: m.Challenge}

		w.Header().Set("Content-Type", "application/json")
		j, _ := json.Marshal(challengeResponse)
		w.Write(j)
	}
}

func eventCallback(w http.ResponseWriter, r *http.Request, m Message) {
	if m.Event.Type == "message" && m.Event.SubType == "" {
		if strings.Contains(m.Event.Text, config.BotName) {
			c.NotifyMention(m)
		} else {
			c.NotifyText(m)
		}
	}
}

func imgurAPI(tag string) string {
	url := "https://api.imgur.com/3/gallery/search/time?q_all=" + url.QueryEscape(tag)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Client-ID "+config.ImgurClient)

	client := &http.Client{}
	client.Timeout = time.Second * 15
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

type sayMessage struct {
	Challenge string `json:"challenge"`
	Channel   string `json:"channel"`
	Message   string `json:"message"`
}

type imgurresp struct {
	Images []imguritem `json:"images"`
	Link   string      `json:"link"`
}

type imguritem struct {
	Link string `json:"link"`
}

type challengeResponse struct {
	Challenge string `json:"challenge"`
}

type configuration struct {
	VerificationToken string `json:"verification_token"`
	BotToken          string `json:"bot_token"`
	ImgurClient       string `json:"imgur_client"`
	BotName           string `json:"bot_name"`
	SayChallenge      string `json:"say_challenge"`
}
