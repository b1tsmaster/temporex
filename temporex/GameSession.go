package main

import "time"

type GameSession struct {
	sessionID string
	playerIDs []string
	createdAt time.Time
	groups    map[string][]string // groupID to playerIDs
}

func (ms *GameSession) AssignPlayerToGroup(playerID, groupID string) {
	if ms.groups == nil {
		ms.groups = make(map[string][]string)
	}
	ms.groups[groupID] = append(ms.groups[groupID], playerID)
}

func ForwardMsg(sessionID string, groupID string, message []byte) {

}
