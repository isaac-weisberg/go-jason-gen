package sample

import gojason "github.com/isaac-weisberg/go-jason"

type accessTokenHaving struct {
	gojason.Decodable

	accessToken string
}
