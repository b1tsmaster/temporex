package main

import (
	"encoding/json"
	"log"
	"os"
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

func MakeMatch() *MatchSession {
	return &MatchSession{}
}

func ClearMatch(session *MatchSession) {
	session.playerIDs = []string{}
	session.groups = make(map[string][]string)
}
