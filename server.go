package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var config configuration

func main() {
	file, err := os.Open("conf.json")

	if err != nil {
		log.Fatal(err)
	}

	_ = json.NewDecoder(file).Decode(&config)

	router := mux.NewRouter()
	router.HandleFunc("/", postMessage).Methods("POST")
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
			NotifyMention(m)
		} else {
			NotifyText(m)
		}
	}
}

func imgurAPI(tag string) string {
	url := "https://api.imgur.com/3/gallery/t/dog/top/day/"
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

func postResponse(channel string, text string) {
	url := "https://slack.com/api/chat.postMessage"
	j, _ := json.Marshal(Reply{Text: text, Channel: channel})
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(j)))
	req.Header.Set("Authorization", "Bearer "+config.BotToken)
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

type imgurresp struct {
	Items []imguritem `json:"items"`
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
}
