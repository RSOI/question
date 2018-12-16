package controller

import (
	"encoding/json"

	"github.com/RSOI/question/model"
	"github.com/RSOI/question/view"
)

// UpdatePATCH remove question
func UpdatePATCH(body []byte) (*model.Question, error) {
	var err error

	var QuestionToUpdate model.Question
	var UpdatedQuestion model.Question
	err = json.Unmarshal(body, &QuestionToUpdate)
	if err != nil {
		return nil, err
	}

	err = view.ValidateUpdateQuestion(QuestionToUpdate)
	if err != nil {
		return nil, err
	}

	UpdatedQuestion, err = QuestionModel.UpdateQuestion(QuestionToUpdate)
	if err != nil {
		return nil, err
	}
	return &UpdatedQuestion, nil
}
