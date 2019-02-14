package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// SlackMessager is a implementation of the messenger interface that posts to slack's api
type slackMessager struct {
	Config configuration
}

func (s slackMessager) postResponse(channel string, text string) error {
	url := "https://slack.com/api/chat.postMessage"
	j, _ := json.Marshal(Reply{Text: text, Channel: channel})
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(j)))
	req.Header.Set("Authorization", "Bearer "+s.Config.BotToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	client.Timeout = time.Second * 15
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Print(string(body))
	return nil
}
