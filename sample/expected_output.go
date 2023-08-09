package sample

import (
	"errors"

	gojason "github.com/isaac-weisberg/go-jason"
	values "github.com/isaac-weisberg/go-jason/values"
)

func makeAddMoneyRequestFromBytes(bytes []byte) (*addMoneyRequest, error) {
	var j = errors.Join
	var e = errors.New

	rootValueAny, err := gojason.Parse(bytes)
	if err != nil {
		return nil, j(e("parsing json into an object tree failed"), err)
	}

	rootObject, err := rootValueAny.AsObject()
	if err != nil {
		return nil, j(e("interpreting root json value as an object failed"), err)
	}

	parsedObject, err := parseAddMoneyRequestFromJsonObject(rootObject)
	if err != nil {
		return nil, j(e("parsing json into the resulting value failed"), err)
	}

	return parsedObject, nil
}

func parseAddMoneyRequestFromJsonObject(rootObject *values.JsonValueObject) (*addMoneyRequest, error) {
	var j = errors.Join
	var e = errors.New

	var stringKeyValues = rootObject.StringKeyedKeyValuesOnly()

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

	valueForAccessTokenKey, exists := stringKeyValues["accessToken"]
	if !exists {
		return nil, j(e("value not found for key 'accessToken'"))
	}
	valueForAccessTokenKeyAsStringValue, err := valueForAccessTokenKey.AsString()
	if err != nil {
		return nil, j(e("parsing value for key 'accessToken failed"), err)
	}
	parsedStringForAccessTokenKey := valueForAccessTokenKeyAsStringValue.String

	var accessTokenHaving = accessTokenHaving{
		accessToken: parsedStringForAccessTokenKey,
	}

	var decodable = gojason.Decodable{}

	var resultingStructAddMoneyRequest = addMoneyRequest{
		Decodable:         decodable,
		amount:            *parsedInt64ForAmountKey,
		accessTokenHaving: accessTokenHaving,
	}

	return &resultingStructAddMoneyRequest, nil
}
