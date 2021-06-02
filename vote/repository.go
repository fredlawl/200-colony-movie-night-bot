package vote

import (
	"database/sql"
	"log"

	"github.com/pkg/errors"

	"github.com/fredlawl/200-colony-movie-night-bot/general"
	"github.com/fredlawl/200-colony-movie-night-bot/suggestion"
)

type Repository struct {
	session *sql.DB
}

type BulkVoteResult struct {
	err  error
	vote Vote
}

func NewRepository(session *sql.DB) *Repository {
	return &Repository{
		session: session,
	}
}

func (context *Repository) BulkSaveVotes(author string, week general.WeekID, votes []Vote) ([]BulkVoteResult, error) {
	emptyBulkResult := []BulkVoteResult{}

	if len(votes) == 0 {
		return emptyBulkResult, nil
	}

	tx, err := context.session.Begin()
	if err != nil {
		return emptyBulkResult, errors.Wrap(err, "")
	}

	truncateStmt, err := context.session.Prepare(`DELETE FROM votes WHERE weekID = ? AND author = ?`)
	if err != nil {
		tx.Rollback()
		return emptyBulkResult, errors.Wrap(err, "")
	}

	_, truncateErr := tx.Stmt(truncateStmt).Exec(week.String(), author)
	if truncateErr != nil {
		tx.Rollback()
		return emptyBulkResult, errors.Wrap(truncateErr, "")
	}

	// TODO: Figure out how to BULK insert w/ prepared statement
	stmt, err := context.session.Prepare(`
		INSERT INTO votes (suggestionID, weekID, author, preference)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		tx.Rollback()
		return emptyBulkResult, errors.Wrap(err, "")
	}

	var hasErrors = false
	var bulkResults = make([]BulkVoteResult, len(votes))
	for i, v := range votes {
		bulkResults[i].vote = v
		_, bulkResults[i].err = tx.Stmt(stmt).Exec(v.SuggestionOrderedID,
			v.WeekID.String(), v.Author, v.Preference)
		bulkResults[i].err = errors.Wrap(bulkResults[i].err, "")
		if !hasErrors {
			hasErrors = bulkResults[i].err != nil
		}
	}

	if hasErrors {
		return bulkResults, tx.Rollback()
	}

	return bulkResults, tx.Commit()
}

func (context *Repository) SuggestionCnt(weekID general.WeekID) int {
	stmt, err := context.session.Prepare("SELECT COUNT(id) FROM suggestions WHERE weekID = ?")
	if err != nil {
		log.Fatalf("[error] %+v\n", errors.Wrap(err, ""))
		return 0
	}

	var cnt int
	queryErr := stmt.QueryRow(weekID.String()).Scan(&cnt)
	if queryErr != nil {
		log.Fatalf("[error] %+v", errors.Wrap(err, ""))
		return 0
	}

	return cnt
}

func (context *Repository) VotesByWeek(weekID general.WeekID) []Vote {
	var votes []Vote = make([]Vote, 0)

	stmt, err := context.session.Prepare(`
	SELECT
		v.suggestionID
		, s.movie
		, v.preference
		, v.author
	FROM votes v
	INNER JOIN suggestions s ON s.id = v.suggestionID
	WHERE weekID = ?`)

	if err != nil {
		return votes
	}

	rows, err := stmt.Query(weekID)
	if err != nil {
		return votes
	}

	var suggestionID int
	var movie string
	var preference uint
	var author string

	for rows.Next() {
		err = rows.Scan(&suggestionID, &movie, &preference, &author)
		if err != nil {
			return votes
		}

		votes = append(votes, Vote{
			SuggestionOrderedID: suggestion.OrderedID(suggestionID),
			Movie:               general.MovieFromString(movie),
			Preference:          preference,
			Author:              author,
		})
	}

	return votes
}
