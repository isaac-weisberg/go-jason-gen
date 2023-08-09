package main

import (
	"errors"
	"strings"
	"unicode"
)

func firstCapitalized(value string) string {
	var copy = value
	if len(value) == 0 {
		return value
	}

	var firstChar = copy[:1]
	var otherChars = copy[1:]

	var builder strings.Builder

	for _, r := range firstChar {
		builder.WriteRune(unicode.ToUpper(r))
	}

	builder.WriteString(otherChars)

	return builder.String()
}

func w(err error, msg string) error {
	return errors.Join(errors.New(msg), err)
}

func e(msg string) error {
	return errors.New(msg)
}
