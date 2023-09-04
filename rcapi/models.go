package rcapi

import _ "encoding/json"

type LoginResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Is2FaEnabled bool   `json:"is_2fa_enabled"`
		Token        string `json:"token"`
		RedirectPage string `json:"redirectPage"`
	} `json:"data"`
	Error string `json:"error"`
}

type GameStartResponse struct {
	Cmd    string `json:"cmd"`
	Cmdval struct {
		UserGameID string `json:"user_game_id"`
		GameNumber int    `json:"game_number"`
		Level      struct {
			Level    int `json:"level"`
			Progress int `json:"progress"`
			Size     int `json:"size"`
		} `json:"level"`
		CoolDown int `json:"cool_down"`
	} `json:"cmdval"`
}

type RCEncryptResponse struct {
	Success bool   `json:"success"`
	Data    string `json:"data"`
	Error   string `json:"error"`
}

type GamesData struct {
	Cmd    string `json:"cmd"`
	Cmdval []struct {
		GameNumber int `json:"game_number"`
		Level      struct {
			Level    int `json:"level"`
			Progress int `json:"progress"`
			Size     int `json:"size"`
		} `json:"level"`
		CoolDown    int `json:"cool_down"`
		CoolDownMax int `json:"cool_down_max"`
	} `json:"cmdval"`
}

type ProfileData struct {
	Cmd    string `json:"cmd"`
	Cmdval struct {
		Auth   bool   `json:"auth"`
		Sessid string `json:"sessid"`
		Userid string `json:"userid"`
	} `json:"cmdval"`
}

type MarketUpdate struct {
	Cmd   string `json:"cmd"`
	Value struct {
		Currency string `json:"currency"`
		ItemType string `json:"item_type"`
		ItemID   string `json:"item_id"`
		Data     struct {
			TradeOffers     string `json:"tradeOffers"`
			AllowedPriceMin int    `json:"allowedPriceMin"`
			AllowedPriceMax int    `json:"allowedPriceMax"`
			List            []struct {
				Price    int `json:"price"`
				Quantity int `json:"quantity"`
			} `json:"list"`
			Rest struct {
				Price    int `json:"price"`
				Quantity int `json:"quantity"`
			} `json:"rest"`
			TotalCount int `json:"totalCount"`
		} `json:"data"`
	} `json:"value"`
}

type Item struct {
	ItemID     string `json:"itemId"`
	ItemType   string `json:"itemType"`
	TotalCount int    `json:"totalCount"`
	Currency   string `json:"currency"`
	TotalPrice int    `json:"totalPrice"`
}