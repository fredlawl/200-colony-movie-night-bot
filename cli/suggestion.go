package main

import (
	"errors"

	"github.com/google/uuid"
)

type SuggestionID string
type SuggestionOrderID uint64

func (SuggestionId SuggestionID) String() string {
	return string(SuggestionId)
}

type Suggestion struct {
	ID     SuggestionID
	WeekID WeekID
	Author string
	Movie  Movie
	Order  SuggestionOrderID
}

func NewSuggestion(weekID WeekID, author string, movie Movie) (*Suggestion, error) {
	if len(movie.Encode()) == 0 {
		return nil, errors.New("movie could not be encoded")
	}

	suggestionID := SuggestionID(uuid.New().String())
	return &Suggestion{
		ID:     suggestionID,
		WeekID: weekID,
		Author: author,
		Movie:  movie,
		Order:  1,
	}, nil
}
