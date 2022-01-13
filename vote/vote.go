package vote

import (
	"github.com/fredlawl/200-colony-movie-night-bot/general"
	"github.com/fredlawl/200-colony-movie-night-bot/suggestion"
)

type ID string

type Vote struct {
	VoteID              ID
	SuggestionOrderedID suggestion.OrderedID
	WeekID              general.WeekID
	Author              string
	Preference          uint
}
