package kensoy

import (
	"regexp"
	"strings"
)

var alphaNumRegex = regexp.MustCompile("[^a-zA-Z0-9\\s]+")

func tokenize(text string) []string {
	lowercase := strings.ToLower(text)

	alphanum := alphaNumRegex.ReplaceAllString(lowercase, "")

	tokens := strings.Split(alphanum, " ")

	return tokens
}
