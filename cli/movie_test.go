package main

import "testing"

func TestGivenMovieNameOfSpecialCharactersThenCharactersReplaced(t *testing.T) {
	movie := MovieFromString("~!@#$%^&*()-=_+,.<>/?;:'\"[{]}\\| ")
	expected := ""
	actual := movie.Encode()

	if actual != expected {
		t.Fail()
	}
}

func TestGivenMovieNameWithArticlesAreThenReplaced(t *testing.T) {
	movie := MovieFromString("the a an and or")

	expected := ""
	actual := movie.Encode()

	if actual != expected {
		t.Fail()
	}
}

func TestGivenMovieNameWithVowelsAreThenReplaced(t *testing.T) {
	movie := MovieFromString("aeiou")

	expected := ""
	actual := movie.Encode()

	if actual != expected {
		t.Fail()
	}
}

func TestGivenMovieNameWithSubsequentRepeatingCharactersAndVowelsAreThenReplaced(t *testing.T) {
	movie := MovieFromString("Winnie the Pooh")

	expected := "wnph"
	actual := movie.Encode()

	if actual != expected {
		t.Fail()
	}
}