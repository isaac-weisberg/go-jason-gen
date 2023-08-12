package sample

import gojason "github.com/isaac-weisberg/go-jason"

// seeemen
type addMoneyRequest struct {
	gojason.Decodable

	accessTokenHaving

	amount     int64
	message    string
	otherStuff []addMoneyRequest

	moneySpent moneySpentRequest
}
