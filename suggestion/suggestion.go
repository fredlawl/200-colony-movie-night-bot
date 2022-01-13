package suggestion

import (
	"errors"

	"github.com/fredlawl/200-colony-movie-night-bot/general"
	"github.com/google/uuid"
)

type ID string
type OrderedID uint64

func (SuggestionId ID) String() string {
	return string(SuggestionId)
}

type Suggestion struct {
	ID     ID
	WeekID general.WeekID
	Author string
	Movie  general.Movie
	Order  OrderedID
}

func NewSuggestion(weekID general.WeekID, author string, movie general.Movie) (*Suggestion, error) {
	if len(movie.Encode()) == 0 {
		return nil, errors.New("movie could not be encoded")
	}

	suggestionID := ID(uuid.New().String())
	return &Suggestion{
		ID:     suggestionID,
		WeekID: weekID,
		Author: author,
		Movie:  movie,
		Order:  1,
	}, nil
}
