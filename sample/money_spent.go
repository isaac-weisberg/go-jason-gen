package sample

import gojason "github.com/isaac-weisberg/go-jason"

type moneySpentRequest struct {
	gojason.Decodable

	spendAmount int64
}
