package main

type VoteID string

type Vote struct {
	VoteID            VoteID
	SuggestionOrderID SuggestionOrderID
	WeekID            WeekID
	Author            string
	Preference        uint
}
