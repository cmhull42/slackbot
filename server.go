package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var token string

func main() {
	file, err := os.Open("conf.json")

	if err != nil {
		log.Fatal(err)
	}
	var config Configuration
	_ = json.NewDecoder(file).Decode(&config)
	token = config.Verification_token

	router := mux.NewRouter()
	router.HandleFunc("/", PostMessage).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func PostMessage(w http.ResponseWriter, r *http.Request) {
	var message Message
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &message)
	challengeResponse := ChallengeResponse{Challenge: message.Challenge}
	if message.Token == token {
		w.Header().Set("Content-Type", "application/json")
		j, _ := json.Marshal(challengeResponse)
		w.Write(j)
	}
}

type Message struct {
	Token      string `json:"token,omitempty"`
	Challenge  string `json:"challenge,omitempty"`
	Event_type string `json:"type,omitempty"`
}

type ChallengeResponse struct {
	Challenge string `json:"challenge"`
}

type Configuration struct {
	Verification_token string
}
