// Code generated by "go-jason-gen"; DO NOT EDIT.

package sample

import (
	"errors"

	gojason "github.com/isaac-weisberg/go-jason"
	values "github.com/isaac-weisberg/go-jason/values"
)

func makeAccessTokenHavingFromJson(bytes []byte) (*accessTokenHaving, error) {
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

	parsedObject, err := parseAccessTokenHavingFromJsonObject(rootObject)
	if err != nil {
		return nil, j(e("parsing json into the resulting value failed"), err)
	}

	return parsedObject, nil
}

func parseAccessTokenHavingFromJsonObject(rootObject *values.JsonValueObject) (*accessTokenHaving, error) {
	var j = errors.Join
	var e = errors.New

	var stringKeyValues = rootObject.StringKeyedKeyValuesOnly()

	valueForAccessTokenKey, exists := stringKeyValues["accessToken"]
	if !exists {
		return nil, j(e("value not found for key 'accessToken'"))
	}
	valueForAccessTokenKeyAsStringValue, err := valueForAccessTokenKey.AsString()
	if err != nil {
		return nil, j(e("interpreting JsonAny as String failed for key 'accessToken'"), err)
	}
	parsedStringForAccessTokenKey := valueForAccessTokenKeyAsStringValue.String

	var decodable = gojason.Decodable{}
	var resultingStructAccessTokenHaving = accessTokenHaving{
		Decodable: decodable,
		accessToken: parsedStringForAccessTokenKey,
	}
	return &resultingStructAccessTokenHaving, nil
}
func makeAddMoneyRequestFromJson(bytes []byte) (*addMoneyRequest, error) {
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

	valueForEmbeddedAccessTokenHaving, err := parseAccessTokenHavingFromJsonObject(rootObject)
	if err != nil {
		return nil, j(e("parsing embedded struct of type 'accessTokenHaving' failed"), err)
	}

	valueForAmountKey, exists := stringKeyValues["amount"]
	if !exists {
		return nil, j(e("value not found for key 'amount'"))
	}
	valueForAmountKeyAsNumberValue, err := valueForAmountKey.AsNumber()
	if err != nil {
		return nil, j(e("interpreting JsonAny as Number failed for key 'amount'"), err)
	}
	parsedInt64ForAmountKey, err := valueForAmountKeyAsNumberValue.ParseInt64()
	if err != nil {
		return nil, j(e("parsing int64 from Number failed for key 'amount'"), err)
	}

	valueForMessageKey, exists := stringKeyValues["message"]
	if !exists {
		return nil, j(e("value not found for key 'message'"))
	}
	valueForMessageKeyAsStringValue, err := valueForMessageKey.AsString()
	if err != nil {
		return nil, j(e("interpreting JsonAny as String failed for key 'message'"), err)
	}
	parsedStringForMessageKey := valueForMessageKeyAsStringValue.String

	var decodable = gojason.Decodable{}
	var resultingStructAddMoneyRequest = addMoneyRequest{
		Decodable: decodable,
		amount: *parsedInt64ForAmountKey,
		message: parsedStringForMessageKey,
		accessTokenHaving: *valueForEmbeddedAccessTokenHaving,
	}
	return &resultingStructAddMoneyRequest, nil
}
