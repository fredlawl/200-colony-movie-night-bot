package vote

import (
	"sort"

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
	Movie               general.Movie
}

type Votes []Vote

type Candidate struct {
	SuggestionOrderedId suggestion.OrderedID
	Movie               general.Movie
	Votes               uint
	Percent             float32
}

type Candidates []*Candidate

type Leaderboard struct {
	CurPass     uint
	TotalPasses uint
	Candidates  []*Candidate
}

type void struct{}

func (vs Votes) Score() Leaderboard {
	candidateLookup := make(map[suggestion.OrderedID]*Candidate)
	disqualifiedCanidatesLookup := make(map[suggestion.OrderedID]void)
	var candidates = make(Candidates, 0)
	var pass uint = 0
	var numberOfPotentialPasses uint = 0
	var empty void

	// Create unique electors
	for _, v := range vs {
		var candidate *Candidate
		candidate, exists := candidateLookup[v.SuggestionOrderedID]
		if !exists {
			candidate = &Candidate{
				SuggestionOrderedId: v.SuggestionOrderedID,
				Movie:               v.Movie,
				Votes:               0,
				Percent:             0.0,
			}

			candidateLookup[v.SuggestionOrderedID] = candidate
			candidates = append(candidates, candidate)
			numberOfPotentialPasses++
		}
	}

	if numberOfPotentialPasses > 0 {
		pass++
	}

	for {
		if pass == 0 {
			break
		}

		// Aggregates
		var totalVotesForPass = 0
		for _, v := range vs {
			_, inDisqualified := disqualifiedCanidatesLookup[v.SuggestionOrderedID]
			if v.Preference < pass || inDisqualified {
				continue
			}

			totalVotesForPass++
		}

		// Tally votes
		for _, v := range vs {
			_, inDisqualified := disqualifiedCanidatesLookup[v.SuggestionOrderedID]
			if v.Preference < pass || inDisqualified {
				continue
			}

			candidate := candidateLookup[v.SuggestionOrderedID]
			candidate.Votes++
			candidate.Percent = float32(candidate.Votes) / float32(totalVotesForPass)
		}

		// Determine max/min for pass
		topCandidate, bottomCandidate := candidates.minMax()

		// We have a majority, or ran out of passes
		if topCandidate.Percent > float32(0.5) || numberOfPotentialPasses-pass <= 0 {
			break
		}

		disqualifiedCanidatesLookup[bottomCandidate.SuggestionOrderedId] = empty

		pass++
	}

	// Place highest percentage first
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Percent > candidates[j].Percent
	})

	return Leaderboard{
		CurPass:     pass,
		TotalPasses: numberOfPotentialPasses,
		Candidates:  candidates,
	}
}

func (as Candidates) minMax() (max *Candidate, min *Candidate) {
	var top *Candidate = as[0]
	var bottom *Candidate = as[0]

	for _, aggregate := range as {
		if aggregate.Votes > bottom.Votes {
			top = aggregate
		}

		if aggregate.Votes < bottom.Votes {
			bottom = aggregate
		}
	}

	return top, bottom
}
