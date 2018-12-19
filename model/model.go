package model

import (
	"time"

	"github.com/jackc/pgx"
)

// Question interface
type Question struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Content        *string   `json:"content"`
	AuthorID       int       `json:"author_id"`
	AuthorNickname string    `json:"author_nickname"`
	HasBest        *bool     `json:"has_best"`
	Created        time.Time `json:"created"`
}

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
