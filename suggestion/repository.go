package suggestion

import (
	"database/sql"
	"log"

	"github.com/pkg/errors"

	"github.com/fredlawl/200-colony-movie-night-bot/general"
)

type Repository struct {
	session *sql.DB
}

func NewRepository(session *sql.DB) *Repository {
	return &Repository{
		session: session,
	}
}

func (context *Repository) Save(s Suggestion) error {
	stmt, err := context.session.Prepare(
		`INSERT INTO suggestions (
			uuid,
			weekID,
			author,
			movie,
			movieHash
		) VALUES (
			?,
			?,
			?,
			?,
			?
		)`)

	if err != nil {
		return errors.Wrap(err, "")
	}

	_, err = stmt.Exec(s.ID.String(), s.WeekID.String(), s.Author,
		s.Movie.String(), s.Movie.Encode())

	return errors.Wrap(err, "")
}

func (context *Repository) AllSuggestions(weekID general.WeekID, callback func(key []byte, suggestion *Suggestion) error) {
	stmt, err := context.session.Prepare("SELECT id, uuid, author, movie FROM suggestions WHERE weekID = ? ORDER BY id ASC")
	if err != nil {
		log.Fatalf("[error] %+v", errors.Wrap(err, ""))
		return
	}

	rows, err := stmt.Query(weekID.String())
	if err != nil {
		log.Fatalf("[error]  %+v", errors.Wrap(err, ""))
		return
	}

	var id int
	var suggestionID string
	var author string
	var movie string

	for rows.Next() {
		err = rows.Scan(&id, &suggestionID, &author, &movie)
		if err != nil {
			return
		}

		err = callback(
			[]byte(suggestionID),
			&Suggestion{
				ID:     ID(suggestionID),
				WeekID: weekID,
				Author: author,
				Movie:  general.MovieFromString(movie),
				Order:  OrderedID(id),
			})

		if err != nil {
			log.Fatalf("[error] %+v", errors.Wrap(err, ""))
			return
		}
	}
}

// GetSuggestionByOrder Given the order id, return the suggestion at that position
func (context *Repository) GetSuggestionByOrder(orderID OrderedID) *Suggestion {
	stmt, err := context.session.Prepare("SELECT id, uuid, weekID, author, movie FROM suggestions WHERE id = ?")
	if err != nil {
		log.Fatalf("[error] %+v", errors.Wrap(err, ""))
		return nil
	}

	row := stmt.QueryRow(orderID)

	var id int
	var suggestionID string
	var weekID string
	var author string
	var movie string

	err = row.Scan(&id, &suggestionID, &weekID, &author, &movie)
	if err != nil {
		return nil
	}

	parsedWeekID, _ := general.WeekIDFromString(weekID)

	return &Suggestion{
		ID:     ID(suggestionID),
		WeekID: *parsedWeekID,
		Author: author,
		Movie:  general.MovieFromString(movie),
		Order:  OrderedID(id),
	}
}

func (context *Repository) Remove(s Suggestion) error {
	stmt, err := context.session.Prepare("DELETE FROM suggestions WHERE uuid = ?")
	if err != nil {
		return errors.Wrap(err, "")
	}

	_, err = stmt.Exec(s.ID.String())
	return errors.Wrap(err, "")
}
