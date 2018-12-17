package controller

import (
	"encoding/json"
	"fmt"

	"github.com/RSOI/question/model"
	"github.com/RSOI/question/utils"
	"github.com/RSOI/question/view"
)

// RemoveDELETE remove question
func RemoveDELETE(body []byte) error {
	var err error

	var QuestionToRemove model.Question
	err = json.Unmarshal(body, &QuestionToRemove)
	if err != nil {
		utils.LOG(fmt.Sprintf("Broken body. Error: %s", err.Error()))
		return err
	}

	f, err := view.ValidateDeleteQuestion(QuestionToRemove)
	if err != nil {
		utils.LOG(fmt.Sprintf("Validation error: %s", err.Error()))
		return err
	}

	utils.LOG(fmt.Sprintf("Removing question by: %s...", f))

	switch f {
	case "id":
		err = QuestionModel.DeleteQuestionByID(QuestionToRemove)
	case "author_id":
		err = QuestionModel.DeleteQuestionByAuthorID(QuestionToRemove)
	}

	if err != nil {
		utils.LOG(fmt.Sprintf("Data error: %s", err.Error()))
	} else {
		utils.LOG("Question removed successfully")
	}

	return err
}
