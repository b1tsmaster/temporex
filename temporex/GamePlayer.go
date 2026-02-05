package main

import (
	"time"

	"github.com/gorilla/websocket"
)

type PlayerConn struct {
	createdAt time.Time
	conn      *websocket.Conn
}

type GamePlayer struct {
	playerID string
	groups   []string
	conn     *PlayerConn
}

func AddPlayerConn(conn *websocket.Conn) *PlayerConn {
	return &PlayerConn{
		createdAt: time.Now(),
		conn:      conn,
	}
}

func AddGamePlayer(playerID string) *GamePlayer {
	return &GamePlayer{
		playerID: playerID,
		groups:   []string{},
	}
}
