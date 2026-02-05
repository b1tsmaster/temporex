package main

import "time"

type MatchSession struct {
	sessionId string
	playerIDs []string
	createdAt time.Time
	groups    map[string][]string // groupID to playerIDs
}
