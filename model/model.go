package model

import (
	"github.com/jackc/pgx"
)

// QService connection holder
type QService struct {
	Conn *pgx.ConnPool
}

// QServiceInterface question methods interface
type QServiceInterface interface {
	AddQuestion(q Question) (Question, error)
	DeleteQuestionByID(q Question) error
	DeleteQuestionByAuthorID(q Question) error
	GetQuestionByID(qID int) (Question, error)
	GetQuestionsByAuthorID(qAuthorID int) ([]Question, error)
	UpdateQuestion(q Question) (Question, error)
	GetUsageStatistic(host string) (ServiceStatus, error)
	LogStat(request []byte, responseStatus int, responseError string)
}
