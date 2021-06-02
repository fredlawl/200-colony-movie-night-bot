package vote

import (
	"testing"

	"github.com/fredlawl/200-colony-movie-night-bot/general"
	"github.com/fredlawl/200-colony-movie-night-bot/suggestion"

	"github.com/stretchr/testify/assert"
)

func TestScoreReturnsEmptyGivenEmptyVotes(t *testing.T) {
	scores := make(Votes, 0)

	actual := scores.Score()

	assert.Equal(t, 0, len(actual.Candidates))
	assert.Equal(t, uint(0), actual.TotalPasses)
	assert.Equal(t, uint(0), actual.CurPass)
}

func TestScoreReturnsUniqueSuggestions(t *testing.T) {
	scores := onePassMajority()

	actual := scores.Score()

	assert.Equal(t, 3, len(actual.Candidates))
}

func TestScoreReturnsCorrectValuesAfterFirstPass(t *testing.T) {
	scores := onePassMajority()

	actual := scores.Score()

	assert.Equal(t, uint(1), actual.CurPass)
	assert.Equal(t, uint(3), actual.TotalPasses)

	assert.Equal(t, uint(4), actual.Candidates[0].Votes)
	assert.GreaterOrEqual(t, float32(0.667), actual.Candidates[0].Percent)

	assert.Equal(t, uint(1), actual.Candidates[1].Votes)
	assert.GreaterOrEqual(t, float32(0.167), actual.Candidates[1].Percent)

	assert.Equal(t, uint(1), actual.Candidates[2].Votes)
	assert.GreaterOrEqual(t, float32(0.167), actual.Candidates[2].Percent)
}

func TestScoreReturnsCorrectValuesAfterTwoPass(t *testing.T) {
	scores := twoPassMajority()

	actual := scores.Score()

	assert.Equal(t, uint(2), actual.CurPass)
	assert.Equal(t, uint(3), actual.TotalPasses)

	assert.Equal(t, uint(4), actual.Candidates[0].Votes)
	assert.GreaterOrEqual(t, float32(0.667), actual.Candidates[0].Percent)

	assert.Equal(t, uint(1), actual.Candidates[1].Votes)
	assert.GreaterOrEqual(t, float32(0.167), actual.Candidates[1].Percent)

	assert.Equal(t, uint(1), actual.Candidates[2].Votes)
	assert.GreaterOrEqual(t, float32(0.167), actual.Candidates[2].Percent)
}

func onePassMajority() Votes {
	return Votes{
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(1),
			Movie:               general.MovieFromString("test"),
			Preference:          uint(1),
			Author:              "liam",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(1),
			Movie:               general.MovieFromString("test"),
			Preference:          uint(1),
			Author:              "oliver",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(1),
			Movie:               general.MovieFromString("test"),
			Preference:          uint(1),
			Author:              "james",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(1),
			Movie:               general.MovieFromString("test"),
			Preference:          uint(1),
			Author:              "sneaky",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(2),
			Movie:               general.MovieFromString("shreck"),
			Preference:          uint(1),
			Author:              "william",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(3),
			Movie:               general.MovieFromString("pooh"),
			Preference:          uint(1),
			Author:              "noah",
		},
	}
}

func twoPassMajority() Votes {
	return Votes{
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(1),
			Movie:               general.MovieFromString("test"),
			Preference:          uint(1),
			Author:              "liam",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(1),
			Movie:               general.MovieFromString("test"),
			Preference:          uint(1),
			Author:              "james",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(2),
			Movie:               general.MovieFromString("shreck"),
			Preference:          uint(1),
			Author:              "william",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(3),
			Movie:               general.MovieFromString("pooh"),
			Preference:          uint(1),
			Author:              "noah",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(3),
			Movie:               general.MovieFromString("pooh"),
			Preference:          uint(1),
			Author:              "sneaky",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(1),
			Movie:               general.MovieFromString("test"),
			Preference:          uint(2),
			Author:              "william",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(2),
			Movie:               general.MovieFromString("shreck"),
			Preference:          uint(2),
			Author:              "liam",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(2),
			Movie:               general.MovieFromString("shreck"),
			Preference:          uint(2),
			Author:              "noah",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(2),
			Movie:               general.MovieFromString("shreck"),
			Preference:          uint(2),
			Author:              "james",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(2),
			Movie:               general.MovieFromString("shreck"),
			Preference:          uint(2),
			Author:              "sneaky",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(3),
			Movie:               general.MovieFromString("pooh"),
			Preference:          uint(2),
			Author:              "oliver",
		},
	}
}

func threePassMajority() Votes {
	return Votes{
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(1),
			Movie:               general.MovieFromString("test"),
			Preference:          uint(1),
			Author:              "liam",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(1),
			Movie:               general.MovieFromString("test"),
			Preference:          uint(1),
			Author:              "james",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(2),
			Movie:               general.MovieFromString("shreck"),
			Preference:          uint(1),
			Author:              "william",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(3),
			Movie:               general.MovieFromString("pooh"),
			Preference:          uint(1),
			Author:              "noah",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(3),
			Movie:               general.MovieFromString("pooh"),
			Preference:          uint(1),
			Author:              "sneaky",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(1),
			Movie:               general.MovieFromString("test"),
			Preference:          uint(2),
			Author:              "william",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(2),
			Movie:               general.MovieFromString("shreck"),
			Preference:          uint(2),
			Author:              "liam",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(2),
			Movie:               general.MovieFromString("shreck"),
			Preference:          uint(2),
			Author:              "noah",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(2),
			Movie:               general.MovieFromString("shreck"),
			Preference:          uint(2),
			Author:              "james",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(2),
			Movie:               general.MovieFromString("shreck"),
			Preference:          uint(2),
			Author:              "sneaky",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(3),
			Movie:               general.MovieFromString("pooh"),
			Preference:          uint(2),
			Author:              "oliver",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(1),
			Movie:               general.MovieFromString("test"),
			Preference:          uint(3),
			Author:              "noah",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(2),
			Movie:               general.MovieFromString("shreck"),
			Preference:          uint(3),
			Author:              "oliver",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(3),
			Movie:               general.MovieFromString("pooh"),
			Preference:          uint(3),
			Author:              "liam",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(3),
			Movie:               general.MovieFromString("pooh"),
			Preference:          uint(3),
			Author:              "william",
		},
		Vote{
			SuggestionOrderedID: suggestion.OrderedID(3),
			Movie:               general.MovieFromString("pooh"),
			Preference:          uint(3),
			Author:              "james",
		},
	}
}
