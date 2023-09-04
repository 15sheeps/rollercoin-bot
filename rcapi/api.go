package rcapi

import (
	"encoding/json"
	"log"
	"rollercoin-bot/constants"
	"bytes"
)

func (u *RCUser) GetToken() {
    data := map[string]string{
    	"isCaptchaRequired": "False", 
    	"keepSigned": "True",
    	"mail": u.Mail,
    	"password": u.Password,
    	"language":"en",
    }

    jsonData, _ := json.Marshal(data)

    body := u.PostRequestRetry(constants.AuthURL, bytes.NewBuffer(jsonData))

	var result LoginResponse

	if err := json.Unmarshal(body, &result); err != nil {
    	log.Printf("GetToken: failed to unmarshal JSON: %s\n", string(body))

		log.Fatalln(err)
	}

	u.token = result.Data.Token
	log.Printf("GetToken: jwt token obtained: %s\n", u.token)
}

func (u *RCUser) CollectDaily() {
    u.PostRequestRetry(constants.CollectDailyURL, nil)
}

func (u *RCUser) encodeGameData(url string, encryptedData string) string {
    data := map[string]string{
        "data": encryptedData,
    }

    jsonData, _ := json.Marshal(data)

	body := u.PostRequestRetry(url, bytes.NewBuffer(jsonData))

    var result RCEncryptResponse

	if err := json.Unmarshal(body, &result); err != nil {
    	log.Printf("encodeGameData: failed to unmarshal JSON: %s\n", string(body))

		log.Fatalln(err)
	}

    return result.Data 
}