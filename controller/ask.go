package controller

import (
	"encoding/json"
	"fmt"

	"github.com/RSOI/question/model"
	"github.com/RSOI/question/utils"
	"github.com/RSOI/question/view"
)

// AskPUT new question
func AskPUT(body []byte) (*model.Question, error) {
	var err error

	var NewQuestion model.Question
	err = json.Unmarshal(body, &NewQuestion)
	if err != nil {
		utils.LOG(fmt.Sprintf("Broken body. Error: %s", err.Error()))
		return nil, err
	}

	err = view.ValidateNewQuestion(NewQuestion)
	if err != nil {
		utils.LOG(fmt.Sprintf("Validation error: %s", err.Error()))
		return nil, err
	}

	NewQuestion, err = QuestionModel.AddQuestion(NewQuestion)
	if err != nil {
		utils.LOG(fmt.Sprintf("Data error: %s", err.Error()))
		return nil, err
	}

	utils.LOG("New question added successfully")
	return &NewQuestion, nil
}
