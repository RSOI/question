package controller

import (
	"github.com/RSOI/question/model"
	"github.com/jackc/pgx"
)

var (
	// QuestionModel interface with methods
	QuestionModel model.QServiceInterface
)

// Response interface
type Response struct {
	Status int         `json:"status"`
	Error  string      `json:"error"`
	Data   interface{} `json:"data"`
}

// InitModel Init model with pgx connection
func Init(db *pgx.ConnPool) {
	QuestionModel = &model.QService{
		Conn: db,
	}
}

func errToResponse(err error) (int, string) {
	var statusCode int
	var statusText string

	if err != nil {
		statusText = err.Error()
	}

	switch err {
	case nil:
		statusCode = 200
	case pgx.ErrNoRows:
		statusText = model.ErrNoResult.Error()
		statusCode = 404
	case model.ErrNoResult:
		statusCode = 404
	case model.ErrNoDataToDelete:
		statusCode = 404
	case model.ErrNoDataToUpdate:
		statusCode = 404
	case model.ErrUnavailable:
		statusCode = 503
	default:
		statusCode = 500
		//statusText = "Server error. Additional information may be contained in server logs."
	}

	return statusCode, statusText
}
