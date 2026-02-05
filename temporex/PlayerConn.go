package main

import (
	"time"

	"github.com/gorilla/websocket"
)

type PlayerConn struct {
	playerID  string
	createdAt time.Time
	conn      *websocket.Conn
}
