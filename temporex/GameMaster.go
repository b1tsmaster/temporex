package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type ServerConfig struct {
	PlayersPerMatch int `json:"PlayersPerMatch"`
}

var Config ServerConfig

func LoadConfig(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &Config)
	if err != nil {
		return err
	}

	log.Printf("配置加载成功: PlayersPerMatch=%d\n", Config.PlayersPerMatch)
	return nil
}

func MakeMatch() *GameSession {
	return &GameSession{}
}

func ClearMatch(session *GameSession) {
	session.playerIDs = []string{}
	session.groups = make(map[string][]string)
}

var sessions map[string]*GameSession = make(map[string]*GameSession)

func JoinOrCreate(sessionID string) *GameSession {
	if session, exists := sessions[sessionID]; exists {
		return session
	}
	newSession := &GameSession{
		sessionID: sessionID,
		playerIDs: []string{},
		createdAt: time.Now(),
		groups:    make(map[string][]string),
	}
	sessions[sessionID] = newSession
	return newSession
}
