package controller

import (
	"fmt"
	"strconv"

	"github.com/RSOI/question/model"
	"github.com/RSOI/question/utils"
)

// QuestionGET get question by id
func QuestionGET(id string) (*model.Question, error) {
	qID, _ := strconv.Atoi(id)

	Question, err := QuestionModel.GetQuestionByID(qID)
	if err != nil {
		utils.LOG(fmt.Sprintf("Data error: %s", err.Error()))
		return nil, err
	}

	utils.LOG("Question was found successfully")
	return &Question, nil
}

// QuestionsGET get questions by author
func QuestionsGET(aid string) ([]model.Question, error) {
	var err error

	qAuthorID, _ := strconv.Atoi(aid)

	data, err := QuestionModel.GetQuestionsByAuthorID(qAuthorID)
	if err != nil {
		utils.LOG(fmt.Sprintf("Data error: %s", err.Error()))
		return nil, err
	}

	utils.LOG("Questions were found successfully")
	return data, nil
}
