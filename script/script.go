package script

import (
	"rollercoin-bot/rcapi"
	"time"
	"log"
)

func StartBot(user *rcapi.RCUser) {
	user.GetToken()
	
	user.DialWsRetry()

	go func() {
		for {
			msg := user.ReadWsMessage()

			user.HandleWsMessage(msg)
		}
	}()

	user.WriteWsMessage([]byte(`{"cmd":"profile_data"}`))

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		<- ticker.C
		user.WriteWsMessage([]byte(`{"cmd":"games_data_request"}`))
		log.Println("send: games_data_request")
	}		
}