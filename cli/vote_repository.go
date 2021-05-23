package main

import "database/sql"

type VoteRepository struct {
	session *sql.DB
}

type BulkVoteResult struct {
	err  error
	vote Vote
}

func NewVoteRepository(session *sql.DB) *VoteRepository {
	// This could result in an error, but that's fine
	session.Exec("PRAGMA foreign_keys = ON")
	return &VoteRepository{
		session: session,
	}
}

func (context *VoteRepository) BulkSaveVotes(votes []Vote) ([]BulkVoteResult, error) {
	emptyBulkResult := []BulkVoteResult{}

	tx, err := context.session.Begin()
	if err != nil {
		return emptyBulkResult, err
	}

	stmt, err := context.session.Prepare(`
		INSERT OR REPLACE INTO votes (suggestionID, weekID, author, preference)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return emptyBulkResult, tx.Rollback()
	}

	var hasErrors = false
	var bulkResults = make([]BulkVoteResult, len(votes))
	for i, v := range votes {
		bulkResults[i].vote = v
		_, bulkResults[i].err = tx.Stmt(stmt).Exec(v.SuggestionOrderID,
			v.WeekID.String(), v.Author, v.Preference)
		if !hasErrors {
			hasErrors = bulkResults[i].err != nil
		}
	}

	if hasErrors {
		return bulkResults, tx.Rollback()
	}

	return bulkResults, tx.Commit()
}
