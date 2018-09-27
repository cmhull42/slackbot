package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var vtoken string
var btoken string
var imgurclient string

func main() {
	file, err := os.Open("conf.json")

	if err != nil {
		log.Fatal(err)
	}
	var config configuration
	_ = json.NewDecoder(file).Decode(&config)
	vtoken = config.VerificationToken
	btoken = config.BotToken
	imgurclient = config.ImgurClient

	router := mux.NewRouter()
	router.HandleFunc("/", postMessage).Methods("POST")
	log.Fatal(http.ListenAndServe(":9803", router))
}

func postMessage(w http.ResponseWriter, r *http.Request) {
	var m message
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &m)

	switch m.EventType {
	case "url_verification":
		urlVerification(w, r, m)
	case "event_callback":
		go eventCallback(w, r, m)
	default:
		log.Print(string(b))
	}
}

func urlVerification(w http.ResponseWriter, r *http.Request, m message) {
	if m.Token == vtoken {
		challengeResponse := challengeResponse{Challenge: m.Challenge}

		w.Header().Set("Content-Type", "application/json")
		j, _ := json.Marshal(challengeResponse)
		w.Write(j)
	}
}

func eventCallback(w http.ResponseWriter, r *http.Request, m message) {
	if m.Event.Type == "message" {
		handleMessage(m)
	}
}

func handleMessage(m message) {
	if strings.Contains(m.Event.Text, "fuck") {
		postResponse(m.Event.Channel, "https://youtu.be/hpigjnKl7nI?t=2s")
	}

	if strings.Contains(m.Event.Text, "@B725HJ0DT") && strings.Contains(m.Event.Text, "pup") {
		var j map[string]imgurresp
		puppies := imgurAPI("dog")
		json.NewDecoder(strings.NewReader(puppies)).Decode(&j)

		rand.Seed(time.Now().Unix())
		i := rand.Intn(len(j["data"].Items))

		postResponse(m.Event.Channel, j["data"].Items[i].Link)
	}
}

func imgurAPI(tag string) string {
	url := "https://api.imgur.com/3/gallery/t/dog/top/day/"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Client-ID "+imgurclient)

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

func postResponse(channel string, text string) {
	url := "https://slack.com/api/chat.postMessage"
	j, _ := json.Marshal(reply{Text: text, Channel: channel})
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(j)))
	req.Header.Set("Authorization", "Bearer "+btoken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	client.Timeout = time.Second * 15
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Print(string(body))
}

type message struct {
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	EventType string `json:"type"`
	TeamID    string `json:"team_id"`
	Event     event  `json:"event"`
}

type event struct {
	Type        string `json:"type"`
	User        string `json:"user"`
	Text        string `json:"text"`
	ClientMsgID string `json:"client_msg_id"`
	Time        string `json:"ts"`
	Channel     string `json:"channel"`
	ChannelType string `json:"channel_type"`
}

type imgurresp struct {
	Items []imguritem `json:"items"`
}

type imguritem struct {
	Link string `json:"link"`
}

type challengeResponse struct {
	Challenge string `json:"challenge"`
}

type reply struct {
	Text    string `json:"text"`
	Channel string `json:"channel"`
}

type configuration struct {
	VerificationToken string `json:"verification_token"`
	BotToken          string `json:"bot_token"`
	ImgurClient       string `json:"imgur_client"`
}
