package rcapi

import (
	"github.com/gorilla/websocket"
)

type RCUser struct {
	GamesRemaining int
	
	Mail string
	Password string
	Proxy string

	userid string
	token string

	WSConn *websocket.Conn
}