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
	var config configuration
	_ = json.NewDecoder(file).Decode(&config)
	token = config.VerificationToken

	router := mux.NewRouter()
	router.HandleFunc("/", postMessage).Methods("POST")
	log.Fatal(http.ListenAndServe(":9803", router))
}

func postMessage(w http.ResponseWriter, r *http.Request) {
	var message message
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &message)
	challengeResponse := challengeResponse{Challenge: message.Challenge}
	if message.Token == token {
		w.Header().Set("Content-Type", "application/json")
		j, _ := json.Marshal(challengeResponse)
		w.Write(j)
	}
}

type message struct {
	Token     string `json:"token,omitempty"`
	Challenge string `json:"challenge,omitempty"`
	EventType string `json:"type,omitempty"`
}

type challengeResponse struct {
	Challenge string `json:"challenge"`
}

type configuration struct {
	VerificationToken string
}
