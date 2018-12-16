package controller

import (
	"encoding/json"

	"github.com/RSOI/question/model"
	"github.com/RSOI/question/view"
)

// AskPUT new question
func AskPUT(body []byte) (*model.Question, error) {
	var err error

	var NewQuestion model.Question
	err = json.Unmarshal(body, &NewQuestion)
	if err != nil {
		return nil, err
	}

	err = view.ValidateNewQuestion(NewQuestion)
	if err != nil {
		return nil, err
	}

	NewQuestion, err = QuestionModel.AddQuestion(NewQuestion)
	if err != nil {
		return nil, err
	}
	return &NewQuestion, nil
}
