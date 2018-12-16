package controller

import (
	"strconv"

	"github.com/RSOI/question/model"
)

// QuestionGET get question by id
func QuestionGET(id string) (*model.Question, error) {
	qID, _ := strconv.Atoi(id)

	Question, err := QuestionModel.GetQuestionByID(qID)
	if err != nil {
		return nil, err
	}

	return &Question, nil
}

// QuestionsGET get questions by author
func QuestionsGET(aid string) ([]model.Question, error) {
	var err error

	qAuthorID, _ := strconv.Atoi(aid)

	data, err := QuestionModel.GetQuestionsByAuthorID(qAuthorID)
	if err != nil {
		return nil, err
	}

	return data, nil
}
