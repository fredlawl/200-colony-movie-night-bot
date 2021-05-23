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
	return &VoteRepository{
		session: session,
	}
}

func (context *VoteRepository) BulkSaveVotes(votes []Vote) ([]BulkVoteResult, error) {
	emptyBulkResult := []BulkVoteResult{}

	context.session.Exec("PRAGMA foreign_keys = ON;")

	tx, err := context.session.Begin()
	if err != nil {
		return emptyBulkResult, err
	}

	stmt, err := context.session.Prepare(`
		INSERT INTO votes (suggestionID, weekID, author, preference)
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

	var txErr error
	if hasErrors {
		txErr = tx.Rollback()
	} else {
		txErr = tx.Commit()
	}

	context.session.Exec("PRAGMA foreign_keys = OFF;")

	return bulkResults, txErr
}
