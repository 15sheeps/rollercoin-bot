package rcapi

import (
	"rollercoin-bot/utils"
	"rollercoin-bot/constants"
	"fmt"

	"time"
)

func (user RCUser) encodeStartGameData(encryptedData string) string {
	url := constants.EncodeStartURL + user.userid
    return user.encodeGameData(url, encryptedData)
}

func (user RCUser) encodeEndGameData(encryptedData string) string {
	url := constants.EncodeEndURL + user.userid
    return user.encodeGameData(url, encryptedData)
}

// добавить обработку ошибок
func (u RCUser) getGameStartCmd(game_number int) string {
	var cmd = fmt.Sprintf(`{"game_number":%d}`, game_number)
	cmd = utils.EncryptCmd(cmd, u.userid)
	cmd = u.encodeStartGameData(cmd)
	return fmt.Sprintf(`{"cmd":"game_start_request","cmdval":"%s"}`, cmd)
}
// добавить обработку ошибок
func (u RCUser) getGameEndCmd(power int, user_game_id string) string {
	t := time.Now().UnixMilli()

	var cmd = fmt.Sprintf(`{"power":%d,"time":%d,"user_game_id":"%s","win_status":3,"token":""}`, power, t, user_game_id)
	cmd = utils.EncryptCmd(cmd, u.userid)
	cmd = u.encodeEndGameData(cmd)

	return fmt.Sprintf(`{"cmd":"game_end_request","cmdval":"%s"}`, cmd)
}