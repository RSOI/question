package controller

import (
	"encoding/json"

	"github.com/RSOI/question/model"
	"github.com/RSOI/question/view"
)

// RemoveDELETE remove question
func RemoveDELETE(body []byte) error {
	var err error

	var QuestionToRemove model.Question
	err = json.Unmarshal(body, &QuestionToRemove)
	if err != nil {
		return err
	}

	f, err := view.ValidateDeleteQuestion(QuestionToRemove)
	if err != nil {
		return err
	}

	switch f {
	case "id":
		err = QuestionModel.DeleteQuestionByID(QuestionToRemove)
	case "author_id":
		err = QuestionModel.DeleteQuestionByAuthorID(QuestionToRemove)
	}

	return err
}
