package main

import "log"

type testMessager struct {
	Config configuration
}

func (t testMessager) postResponse(channel string, text string) error {
	log.Print("channel: " + channel + " - " + text)
	return nil
}
