package controller

import (
	"encoding/json"
	"fmt"

	"github.com/RSOI/question/model"
	"github.com/RSOI/question/utils"
	"github.com/RSOI/question/view"
)

// UpdatePATCH remove question
func UpdatePATCH(body []byte) (*model.Question, error) {
	var err error

	var QuestionToUpdate model.Question
	var UpdatedQuestion model.Question
	err = json.Unmarshal(body, &QuestionToUpdate)
	if err != nil {
		utils.LOG(fmt.Sprintf("Broken body. Error: %s", err.Error()))
		return nil, err
	}

	err = view.ValidateUpdateQuestion(QuestionToUpdate)
	if err != nil {
		utils.LOG(fmt.Sprintf("Validation error: %s", err.Error()))
		return nil, err
	}

	UpdatedQuestion, err = QuestionModel.UpdateQuestion(QuestionToUpdate)
	if err != nil {
		utils.LOG(fmt.Sprintf("Data error: %s", err.Error()))
		return nil, err
	}

	utils.LOG("Question updated successfully")
	return &UpdatedQuestion, nil
}
