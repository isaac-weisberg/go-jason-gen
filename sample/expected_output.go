package sample

import (
	"errors"

	gojason "github.com/isaac-weisberg/go-jason"
)

func newAddMoneyRequest(bytes []byte) (*addMoneyRequest, error) {
	var j = errors.Join
	var e = errors.New

	rootValue, err := gojason.Parse(bytes)
	if err != nil {
		return nil, j(e("parsing json failed"), err)
	}

	rootObject, err := rootValue.AsObject()
	if err != nil {
		return nil, j(e("interpreting root value as an object failed"), err)
	}

	var stringKeyValues = rootObject.StringKeyedKeyValuesOnly()

	var exists bool
	valueForAmountKey, exists := stringKeyValues["amount"]
	if !exists {
		return nil, j(e("value not found for key 'amount'"))
	}
	valueForAmountKeyAsNumberValue, err := valueForAmountKey.AsNumber()
	if err != nil {
		return nil, j(e("parsing value for key 'amount' failed"), err)
	}

	parsedInt64ForAmountKey, err := valueForAmountKeyAsNumberValue.ParseInt64()
	if err != nil {
		return nil, j(e("parsing value for key 'amount' failed"), err)
	}

	var resultingStructAddMoneyRequest = addMoneyRequest{
		amount: *parsedInt64ForAmountKey,
	}

	return &resultingStructAddMoneyRequest, nil
}