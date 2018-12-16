package model

import (
	"time"

	"github.com/RSOI/question/ui"
)

// Question interface
type Question struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Content  *string   `json:"content"`
	AuthorID int       `json:"author_id"`
	HasBest  *bool     `json:"has_best"`
	Created  time.Time `json:"created"`
}

// AddQuestion add new question
func (service *QService) AddQuestion(q Question) (Question, error) {
	var err error

	row := service.Conn.QueryRow(`
		INSERT INTO question.question
			(title, content, author_id) VALUES ($1, $2, $3)
			RETURNING id, created, has_best
	`, q.Title, q.Content, q.AuthorID)

	err = row.Scan(&q.ID, &q.Created, &q.HasBest)
	return q, err
}

// DeleteQuestionByID delete question by id
func (service *QService) DeleteQuestionByID(q Question) error {
	res, err := service.Conn.Exec(`DELETE FROM question.question WHERE id = $1`, q.ID)
	if err == nil && res.RowsAffected() != 1 {
		err = ui.ErrNoDataToDelete
	}
	return err
}

// DeleteQuestionByAuthorID delete question by author id
func (service *QService) DeleteQuestionByAuthorID(q Question) error {
	res, err := service.Conn.Exec(`DELETE FROM question.question WHERE author_id = $1`, q.AuthorID)
	if err == nil && res.RowsAffected() != 1 {
		err = nil
	}
	return err
}

// GetQuestionByID get question data by it's id
func (service *QService) GetQuestionByID(qID int) (Question, error) {
	var err error
	var q Question

	row := service.Conn.QueryRow(`SELECT * FROM question.question WHERE id = $1`, qID)

	err = row.Scan(
		&q.ID,
		&q.Title,
		&q.Content,
		&q.AuthorID,
		&q.HasBest,
		&q.Created)

	return q, err
}

// GetQuestionsByAuthorID get question data by it's id
func (service *QService) GetQuestionsByAuthorID(qAuthorID int) ([]Question, error) {
	var err error
	q := make([]Question, 0)

	rows, err := service.Conn.Query(`SELECT * FROM question.question WHERE author_id = $1 ORDER BY id ASC`, qAuthorID)
	if err != nil {
		return q, err
	}

	for rows.Next() {
		var tq Question
		err = rows.Scan(
			&tq.ID,
			&tq.Title,
			&tq.Content,
			&tq.AuthorID,
			&tq.HasBest,
			&tq.Created)

		if err != nil {
			return q, err
		}

		q = append(q, tq)
	}

	return q, err
}

// UpdateQuestion update question with new data
func (service *QService) UpdateQuestion(q Question) (Question, error) {
	currentQuestionData, err := service.GetQuestionByID(q.ID)
	if err != nil {
		return q, err
	}

	if (q.Content != nil && q.Content != currentQuestionData.Content) ||
		(q.HasBest != nil && q.HasBest != currentQuestionData.HasBest) {

		var content string
		if q.Content != nil {
			content = *q.Content
		} else {
			content = *currentQuestionData.Content
		}

		var best bool
		if q.HasBest != nil {
			best = *q.HasBest
		} else {
			best = *currentQuestionData.HasBest
		}

		res, err := service.Conn.Exec(`
			UPDATE question.question SET content = $1, has_best = $2 WHERE id = $3`,
			content, best, q.ID)
		if err == nil && res.RowsAffected() != 1 {
			err = ui.ErrNoDataToUpdate
		}
	}

	return q, err
}
