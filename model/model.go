package model

import (
	"errors"

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

var (
	// ErrNoResult - no data found
	ErrNoResult = errors.New("no data found")
	// ErrNoDataToDelete - no data found to delete"
	ErrNoDataToDelete = errors.New("no data found to delete")
	// ErrUnavailable - database is unavailable
	ErrUnavailable = errors.New("database is unavailable")
)
