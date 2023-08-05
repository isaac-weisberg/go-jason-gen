package sample

import (
	"errors"

	gojason "github.com/isaac-weisberg/go-jason"
)

var a = 3

func newAddMoneyRequest(jsonString string) (*addMoneyRequest, error) {
	var rootValue, err = gojason.Parse(jsonString)
	if err != nil {
		return nil, errors.Join(errors.New("parsing json failed"), err)
	}

	rootObject, err := rootValue.AsObject()
	if err != nil {
		return nil, errors.Join(errors.New("interpreting root value as an object failed"), err)
	}

}
