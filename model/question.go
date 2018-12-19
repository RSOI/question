package model

import (
	"fmt"

	"github.com/RSOI/question/ui"
	"github.com/RSOI/question/utils"
)

// AddQuestion add new question
func (service *QService) AddQuestion(q Question) (Question, error) {
	var err error

	utils.LOG("Accessing database...")
	row := service.Conn.QueryRow(`
		INSERT INTO question.question
			(title, content, author_id, author_nickname) VALUES ($1, $2, $3, $4)
			RETURNING id, created, has_best
	`, q.Title, q.Content, q.AuthorID, q.AuthorNickname)

	err = row.Scan(&q.ID, &q.Created, &q.HasBest)
	return q, err
}

// DeleteQuestionByID delete question by id
func (service *QService) DeleteQuestionByID(q Question) error {
	utils.LOG("Accessing database...")
	res, err := service.Conn.Exec(`DELETE FROM question.question WHERE id = $1`, q.ID)
	if err == nil && res.RowsAffected() != 1 {
		err = ui.ErrNoDataToDelete
	}
	return err
}

// DeleteQuestionByAuthorID delete question by author id
func (service *QService) DeleteQuestionByAuthorID(q Question) error {
	utils.LOG("Accessing database...")
	res, err := service.Conn.Exec(`DELETE FROM question.question WHERE author_id = $1`, q.AuthorID)
	if err == nil && res.RowsAffected() != 1 {
		utils.LOG(fmt.Sprintf("Author with id '%d' has no questions", q.AuthorID))
		err = nil
	}
	return err
}

// GetQuestionByID get question data by it's id
func (service *QService) GetQuestionByID(qID int) (Question, error) {
	var err error
	var q Question

	utils.LOG("Accessing database...")
	utils.LOG(fmt.Sprintf("SELECT * FROM question.question WHERE id = %d", qID))
	row := service.Conn.QueryRow(`SELECT * FROM question.question WHERE id = $1`, qID)

	err = row.Scan(
		&q.ID,
		&q.Title,
		&q.Content,
		&q.AuthorID,
		&q.AuthorNickname,
		&q.HasBest,
		&q.Created)

	return q, err
}

// GetQuestionsByAuthorID get question data by it's id
func (service *QService) GetQuestionsByAuthorID(qAuthorID int) ([]Question, error) {
	var err error
	q := make([]Question, 0)

	utils.LOG("Accessing database...")
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
			&tq.AuthorNickname,
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

		utils.LOG("Accessing database...")
		utils.LOG(fmt.Sprintf("PDATE question.question SET content = %s, has_best = %t WHERE id = %d", content, best, q.ID))
		res, err := service.Conn.Exec(`
			UPDATE question.question SET content = $1, has_best = $2 WHERE id = $3`,
			content, best, q.ID)
		if err == nil && res.RowsAffected() != 1 {
			err = ui.ErrNoDataToUpdate
		}
	}

	return q, err
}
