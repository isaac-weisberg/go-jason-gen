package sample

import (
	"errors"
	"fmt"

	gojason "github.com/isaac-weisberg/go-jason"
	parser "github.com/isaac-weisberg/go-jason/parser"
	values "github.com/isaac-weisberg/go-jason/values"
)

func ExpectedMakeAddMoneyRequestFromBytes(bytes []byte) (*addMoneyRequest, error) {
	var j = errors.Join
	var e = errors.New

	rootValueAny, err := parser.Parse(bytes)
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

func ExpectedParseAccessTokenHavingFromJsonObject(rootObject *values.JsonValueObject) (*accessTokenHaving, error) {
	var j = errors.Join
	var e = errors.New

	var stringKeyValues = rootObject.StringKeyedKeyValuesOnly()

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

	return &accessTokenHaving, nil
}

func ExpectedParseAddMoneyRequestFromJsonObject(rootObject *values.JsonValueObject) (*addMoneyRequest, error) {
	var j = errors.Join
	var e = errors.New
	var s = fmt.Sprintf

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

	valueForMessageKey, exists := stringKeyValues["message"]
	if !exists {
		return nil, j(e("value not found for key 'message'"))
	}

	valueForMessageKeyAsStringValue, err := valueForMessageKey.AsString()
	if err != nil {
		return nil, j(e("parsing value for key 'message' failed"), err)
	}

	messageResultingValue := valueForMessageKeyAsStringValue.String

	accessTokenHaving, err := ExpectedParseAccessTokenHavingFromJsonObject(rootObject)
	if err != nil {
		return nil, j(e("parsing embedded struct of type accessTokenHaving failed"), err)
	}

	// valueForMoneySpentKey, exists := stringKeyValues["moneySpent"]
	// if !exists {
	// 	return nil, j(e("value not found for key 'moneySpent'"))
	// }
	// valueForMoneySpentKeyAsObjectValue, err := valueForMoneySpentKey.AsObject()
	// if err != nil {
	// 	return nil, j(e("interpreting JsonAny as Object failed for key moneySpent"), err)
	// }
	// parsedValueForMoneySpentKey, err := parseMoneySpentFromJsonObject(valueForMoneySpentKeyAsObjectValue)
	// if err != nil {
	// 	return nil, j(e("parsing 'moneySpentRequest' from 'Object' failed for key 'moneySpent'"))
	// }

	valueForOtherStuffKey, exists := stringKeyValues["otherStuff"]
	if !exists {
		return nil, j(e("value not found for key 'amount'"))
	}
	valueForOtherStuffKeyAsArrayValue, err := valueForOtherStuffKey.AsArray()
	if err != nil {
		return nil, j(e("interpreting JsonAny as String failed for key 'otherStuff'"))
	}
	resultingArrayForOtherStuffKey := make([]addMoneyRequest, 0, len(valueForOtherStuffKeyAsArrayValue.Values))
	for index, element := range valueForOtherStuffKeyAsArrayValue.Values {
		elementAsObject, err := element.AsObject()
		if err != nil {
			return nil, j(e(s("attempted to interpret value at index '%v' of array for key 'otherStuff' as object, but failed", index)), err)
		}

		parsedValue, err := parseAddMoneyRequestFromJsonObject(elementAsObject)
		if err != nil {
			return nil, j(e(s("failed to parse element at index '%v' of array for key 'otherStuff'", index)), err)
		}

		resultingArrayForOtherStuffKey = append(resultingArrayForOtherStuffKey, *parsedValue)
	}

	var decodable = gojason.Decodable{}

	var resultingStructAddMoneyRequest = addMoneyRequest{
		Decodable:         decodable,
		amount:            *parsedInt64ForAmountKey,
		accessTokenHaving: *accessTokenHaving,
		message:           messageResultingValue,
		otherStuff:        resultingArrayForOtherStuffKey,
		// moneySpent:        *moneySpent,
	}

	return &resultingStructAddMoneyRequest, nil
}
