package vote

import (
	"github.com/fredlawl/200-colony-movie-night-bot/cli/general"
	"github.com/fredlawl/200-colony-movie-night-bot/cli/suggestion"
)

type ID string

type Vote struct {
	VoteID              ID
	SuggestionOrderedID suggestion.OrderedID
	WeekID              general.WeekID
	Author              string
	Preference          uint
}
