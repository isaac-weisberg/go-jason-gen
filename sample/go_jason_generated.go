// Code generated by "go-jason-gen"; DO NOT EDIT.

package sample

import (
	"errors"
	"fmt"

	gojason "github.com/isaac-weisberg/go-jason"
	parser "github.com/isaac-weisberg/go-jason/parser"
	values "github.com/isaac-weisberg/go-jason/values"
)

// just in case, if variable is unused, we mention it as a param to this no-op func - and then, it's suddenly very well used :)
func UNUSED(arg any) {}

func makeAccessTokenHavingFromJson(bytes []byte) (*accessTokenHaving, error) {
	var j = errors.Join
	var e = errors.New
	UNUSED(fmt.Sprintf)

	rootValueAny, err := parser.Parse(bytes)
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
	UNUSED(stringKeyValues)

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
	UNUSED(fmt.Sprintf)

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

func parseAddMoneyRequestFromJsonObject(rootObject *values.JsonValueObject) (*addMoneyRequest, error) {
	var j = errors.Join
	var e = errors.New

	var stringKeyValues = rootObject.StringKeyedKeyValuesOnly()
	UNUSED(stringKeyValues)

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

	valueForOtherStuffKey, exists := stringKeyValues["otherStuff"]
	if !exists {
		return nil, j(e("value not found for key 'otherStuff'"))
	}
	valueForOtherStuffKeyAsArrayValue, err := valueForOtherStuffKey.AsArray()
	if err != nil {
		return nil, j(e("interpreting JsonAny as Array failed for key 'otherStuff'"))
	}
	resultingArrayForOtherStuffKey := make([]int64, 0, len(valueForOtherStuffKeyAsArrayValue.Values))
	for index, element := range valueForOtherStuffKeyAsArrayValue.Values {
		elementAsNumber, err := element.AsNumber()
		if err != nil {
			return nil, j(e(fmt.Sprintf("attempted to interpret value at index '%v' of array for key 'otherStuff' as Number, but failed", index)), err)
		}
		parsedInt64, err := elementAsNumber.ParseInt64()
		if err != nil {
			return nil, j(e(fmt.Sprintf("parsing int64 from Number failed for element at index '%v' of array for key 'otherStuff'", index)), err)
		}
		resultingArrayForOtherStuffKey = append(resultingArrayForOtherStuffKey, parsedInt64)
	}

	valueForMoneySpentKey, exists := stringKeyValues["moneySpent"]
	if !exists {
		return nil, j(e("value not found for key 'moneySpent'"))
	}
	valueForMoneySpentKeyAsObjectValue, err := valueForMoneySpentKey.AsObject()
	if err != nil {
		return nil, j(e("interpreting JsonAny as Object failed for key 'moneySpent'"), err)
	}
	parsedValueForMoneySpentKey, err := parseMoneySpentRequestFromJsonObject(valueForMoneySpentKeyAsObjectValue)
	if err != nil {
		return nil, j(e("parsing 'moneySpentRequest' from 'Object' failed for key 'moneySpent'"))
	}

	var decodable = gojason.Decodable{}
	var resultingStructAddMoneyRequest = addMoneyRequest{
		Decodable: decodable,
		amount: parsedInt64ForAmountKey,
		message: parsedStringForMessageKey,
		otherStuff: resultingArrayForOtherStuffKey,
		moneySpent: *parsedValueForMoneySpentKey,
		accessTokenHaving: *valueForEmbeddedAccessTokenHaving,
	}
	return &resultingStructAddMoneyRequest, nil
}

func makeMoneySpentRequestFromJson(bytes []byte) (*moneySpentRequest, error) {
	var j = errors.Join
	var e = errors.New
	UNUSED(fmt.Sprintf)

	rootValueAny, err := parser.Parse(bytes)
	if err != nil {
		return nil, j(e("parsing json into an object tree failed"), err)
	}

	rootObject, err := rootValueAny.AsObject()
	if err != nil {
		return nil, j(e("interpreting root json value as an object failed"), err)
	}

	parsedObject, err := parseMoneySpentRequestFromJsonObject(rootObject)
	if err != nil {
		return nil, j(e("parsing json into the resulting value failed"), err)
	}

	return parsedObject, nil
}

func parseMoneySpentRequestFromJsonObject(rootObject *values.JsonValueObject) (*moneySpentRequest, error) {
	var j = errors.Join
	var e = errors.New

	var stringKeyValues = rootObject.StringKeyedKeyValuesOnly()
	UNUSED(stringKeyValues)

	valueForSpendAmountKey, exists := stringKeyValues["spendAmount"]
	if !exists {
		return nil, j(e("value not found for key 'spendAmount'"))
	}
	valueForSpendAmountKeyAsNumberValue, err := valueForSpendAmountKey.AsNumber()
	if err != nil {
		return nil, j(e("interpreting JsonAny as Number failed for key 'spendAmount'"), err)
	}
	parsedInt64ForSpendAmountKey, err := valueForSpendAmountKeyAsNumberValue.ParseInt64()
	if err != nil {
		return nil, j(e("parsing int64 from Number failed for key 'spendAmount'"), err)
	}

	var decodable = gojason.Decodable{}
	var resultingStructMoneySpentRequest = moneySpentRequest{
		Decodable: decodable,
		spendAmount: parsedInt64ForSpendAmountKey,
	}
	return &resultingStructMoneySpentRequest, nil
}

