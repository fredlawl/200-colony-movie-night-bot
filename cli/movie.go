package main

import (
	"regexp"
	"strings"

	"golang.org/x/text/unicode/norm"
)

type Movie string

func MovieFromString(s string) Movie {
	return Movie(norm.NFD.String(s))
}

func (m Movie) String() string {
	return string(m)
}

func (m Movie) Encode() string {
	pattern := regexp.MustCompile(`[~!@#\$%\^&\*\(\)-=_\+,\.<>\/\?;:'\\"\[{\]}\\\|\s]`)
	encoded :=
		removeVowlesSpacesSuccessiveCharacters(
			removeArticles(
				pattern.ReplaceAllString(
					strings.ToLower(m.String()), " ")))
	return encoded
}

func removeVowlesSpacesSuccessiveCharacters(s string) string {
	var buf strings.Builder
	var last rune
	for _, r := range s {
		if r == last || r == 'a' || r == 'e' || r == 'i' || r == 'o' || r == 'u' || r == ' ' {
			continue
		}

		buf.WriteRune(r)
		last = r
	}
	return buf.String()
}

// a/an/the are the primary articles, but this also counts "or" and "and".
func removeArticles(s string) string {
	articles := map[string]bool{
		"a":   true,
		"an":  true,
		"the": true,
		"and": true,
		"or":  true,
	}

	var buf strings.Builder
	var wordBuf strings.Builder
	var word string

	for _, r := range s {
		if r == ' ' {
			word = wordBuf.String()

			wordBuf.Reset()

			if articles[word] {
				continue
			}

			buf.WriteString(word)
			buf.WriteRune(r)
			continue
		}

		wordBuf.WriteRune(r)
	}

	word = wordBuf.String()
	if !articles[word] {
		buf.WriteString(word)
	}

	return buf.String()
}
