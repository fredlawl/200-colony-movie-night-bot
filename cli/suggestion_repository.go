package main

import "database/sql"

type SuggestionRepository struct {
	session *sql.DB
}

func NewSuggestionRepository(session *sql.DB) *SuggestionRepository {
	return &SuggestionRepository{
		session: session,
	}
}

func (context *SuggestionRepository) Save(s Suggestion) error {
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
		return err
	}

	_, err = stmt.Exec(s.ID.String(), s.WeekID.String(), s.Author,
		s.Movie.String(), s.Movie.Encode())

	return err
}

func (context *SuggestionRepository) AllSuggestions(weekID WeekID, callback func(key []byte, suggestion *Suggestion) error) {
	stmt, err := context.session.Prepare("SELECT id, uuid, author, movie FROM suggestions WHERE weekID = ? ORDER BY id ASC")
	if err != nil {
		return
	}

	rows, err := stmt.Query(weekID.String())
	if err != nil {
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
				ID:     SuggestionID(suggestionID),
				WeekID: weekID,
				Author: author,
				Movie:  MovieFromString(movie),
				Order:  SuggestionOrderID(id),
			})

		if err != nil {
			return
		}
	}
}

// GetSuggestionByOrder Given the order id, return the suggestion at that position
func (context *SuggestionRepository) GetSuggestionByOrder(orderID SuggestionOrderID) *Suggestion {
	stmt, err := context.session.Prepare("SELECT id, uuid, weekID, author, movie FROM suggestions WHERE id = ?")
	if err != nil {
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

	parsedWeekID, _ := WeekIDFromString(weekID)

	return &Suggestion{
		ID:     SuggestionID(suggestionID),
		WeekID: *parsedWeekID,
		Author: author,
		Movie:  MovieFromString(movie),
		Order:  SuggestionOrderID(id),
	}
}

func (context *SuggestionRepository) Remove(s Suggestion) error {
	context.session.Exec("PRAGMA foreign_keys = ON;")
	stmt, err := context.session.Prepare("DELETE FROM suggestions WHERE uuid = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(s.ID.String())

	context.session.Exec("PRAGMA foreign_keys = OFF;")
	return err
}
