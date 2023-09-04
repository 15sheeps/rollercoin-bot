package rcapi

import (
	"rollercoin-bot/constants"
	"time"
	"encoding/json"
	"log"
)

func (user *RCUser) handleGamesDataResponse(resp GamesData) {
	for i := 0; i < len(resp.Cmdval); i++ {
		if (resp.Cmdval[i].CoolDown == 0) {
			cmd := user.getGameStartCmd(resp.Cmdval[i].GameNumber)

			user.WriteWsMessage([]byte(cmd))
			log.Printf("handleWSMessage: send %s.\n", cmd)

			time.Sleep(200 * time.Millisecond)
		}
	}
}

func (user *RCUser) handleGameStartResponse(resp GameStartResponse) {
	var user_game_id = resp.Cmdval.UserGameID
	var level = resp.Cmdval.Level.Level
	var game_number = resp.Cmdval.GameNumber

	var power = constants.RewardsMap[game_number][level]
	// imba
	power = constants.RewardsMap[game_number][10]

	cmd := user.getGameEndCmd(power, user_game_id)

	user.WriteWsMessage([]byte(cmd))
	log.Printf("handleWSMessage: send %s.\n", cmd)
}

func (user *RCUser) HandleWsMessage(data []byte) {
    var cmd struct {
        Cmd string `json:"cmd"`
    }

    if err := json.Unmarshal(data, &cmd); err != nil {
    	log.Printf("handleWSMessage: failed to unmarshal JSON: %s\n", string(data))
    	return
    }

	switch cmd.Cmd {
	case "games_data_response":
		log.Println("handleWSMessage: received games_data_response.")
		var result GamesData
		json.Unmarshal(data, &result)

		if user.GamesRemaining <= 0 {
			return
		}

		user.handleGamesDataResponse(result)
	case "profile":
		log.Println("handleWSMessage: received profile.")
		var result ProfileData
		json.Unmarshal(data, &result)

		user.userid = result.Cmdval.Userid
	case "game_start_response":
		log.Println("handleWSMessage: received game_start_response.")
		var result GameStartResponse
		json.Unmarshal(data, &result)

		user.handleGameStartResponse(result)
	case "game_finished_accepted":
		log.Println("handleWSMessage: received game_finished_accepted.")
		user.GamesRemaining--
	}
}