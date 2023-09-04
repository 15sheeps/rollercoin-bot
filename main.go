package main

import (
	"rollercoin-bot/script"
	"rollercoin-bot/rcapi"
	"time"
	"flag"
)

func main() {
	mail := flag.String("mail", "none", "rollercoin user email")
	password := flag.String("password", "none", "rollercoin user password")
	proxy := flag.String("proxy", "socks5://127.0.0.1:9060", "socks5 proxy")
	games := flag.Int("games", 200, "games to play")

	flag.Parse()

	user := rcapi.RCUser{
		Mail: *mail, 
		Password: *password,
		Proxy: *proxy,
		GamesRemaining: *games,
	}

	go script.StartBot(&user)

	time.Sleep(9 * time.Hour)
}