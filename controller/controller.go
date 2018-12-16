package controller

import (
	"github.com/RSOI/question/model"
	"github.com/jackc/pgx"
)

var (
	// QuestionModel interface with methods
	QuestionModel model.QServiceInterface
)

// Init Init model with pgx connection
func Init(db *pgx.ConnPool) {
	QuestionModel = &model.QService{
		Conn: db,
	}
}
