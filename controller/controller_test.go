package controller

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/RSOI/question/model"
	"github.com/RSOI/question/ui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedQService struct {
	mock.Mock
}

func getMock() *MockedQService {
	QuestionModel = &MockedQService{}
	return QuestionModel.(*MockedQService)
}

func (s *MockedQService) AddQuestion(q model.Question) (model.Question, error) {
	args := s.Mock.Called(q)
	return args.Get(0).(model.Question), args.Error(1)
}
func (s *MockedQService) DeleteQuestionByID(q model.Question) error {
	args := s.Mock.Called(q)
	return args.Error(0)
}
func (s *MockedQService) DeleteQuestionByAuthorID(q model.Question) error {
	args := s.Mock.Called(q)
	return args.Error(0)
}
func (s *MockedQService) GetQuestionByID(qID int) (model.Question, error) {
	args := s.Mock.Called(qID)
	return args.Get(0).(model.Question), args.Error(1)
}
func (s *MockedQService) GetQuestionsByAuthorID(qAuthorID int) ([]model.Question, error) {
	args := s.Mock.Called(qAuthorID)
	return args.Get(0).([]model.Question), args.Error(1)
}
func (s *MockedQService) UpdateQuestion(q model.Question) (model.Question, error) {
	args := s.Mock.Called(q)
	return args.Get(0).(model.Question), args.Error(1)
}
func (s *MockedQService) GetUsageStatistic(host string) (model.ServiceStatus, error) {
	args := s.Mock.Called(host)
	return args.Get(0).(model.ServiceStatus), args.Error(1)
}
func (s *MockedQService) LogStat(request []byte, responseStatus int, responseError string) {
	// nothing interesting here, just store data without affecting main thread
}

var (
	defaultQuestionContent        = "My Content"
	defaultQuestionHasBest        = false
	updatedQuestionHasBest        = true
	defaultQuestionCreatedTime, _ = time.Parse("2006-01-02T15:04:05", time.Now().String())
	defaultQuestion               = model.Question{
		AuthorID: 1,
		Title:    "My Title",
		Content:  &defaultQuestionContent,
	}
	createdQuestion = model.Question{
		ID:       1,
		AuthorID: 1,
		Title:    "My Title",
		Content:  &defaultQuestionContent,
		HasBest:  &defaultQuestionHasBest,
		Created:  defaultQuestionCreatedTime,
	}
	updatedQuestion = model.Question{
		ID:       1,
		AuthorID: 1,
		Title:    "My Title",
		Content:  &defaultQuestionContent,
		HasBest:  &updatedQuestionHasBest,
		Created:  defaultQuestionCreatedTime,
	}
	questionToRemoveID = model.Question{
		ID: 1,
	}
	questionToRemoveAuthorID = model.Question{
		AuthorID: 1,
	}
)

/*
********************************************************************
TESTS FOR ASK ******************************************************
********************************************************************
*/

func TestAskAddCorrectData(t *testing.T) {
	body, _ := json.Marshal(&defaultQuestion)

	cMock := getMock()
	cMock.On("AddQuestion", defaultQuestion).Return(createdQuestion, nil)

	data, err := AskPUT(body)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)

		assert.Equal(t, createdQuestion.ID, data.ID)
		assert.Equal(t, createdQuestion.Title, data.Title)
		assert.Equal(t, *createdQuestion.Content, *data.Content)
		assert.Equal(t, createdQuestion.AuthorID, data.AuthorID)
		assert.Equal(t, *createdQuestion.HasBest, *data.HasBest)
		assert.Equal(t, createdQuestion.Created, data.Created)
	}
}

func TestAskAddMissedField(t *testing.T) {
	body := []byte("{\"author_id\": 1, \"content\": \"My Content\"}")

	data, err := AskPUT(body)
	assert.Equal(t, ui.ErrFieldsRequired, err)
	assert.Nil(t, data)
}

func TestAskAddBrokenBody(t *testing.T) {
	body := []byte("{\"author_id\": 1, \"content\": \"My Content}")

	data, err := AskPUT(body)
	assert.NotNil(t, err)
	assert.Nil(t, data)
}

/*
********************************************************************
TESTS FOR QUESTION ID **********************************************
********************************************************************
*/

func TestQuestionGetByIDCorrectData(t *testing.T) {
	cMock := getMock()
	cMock.On("GetQuestionByID", 1).Return(createdQuestion, nil)

	data, err := QuestionGET("1")
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)

		assert.Equal(t, createdQuestion.ID, data.ID)
		assert.Equal(t, createdQuestion.Title, data.Title)
		assert.Equal(t, *createdQuestion.Content, *data.Content)
		assert.Equal(t, createdQuestion.AuthorID, data.AuthorID)
		assert.Equal(t, *createdQuestion.HasBest, *data.HasBest)
		assert.Equal(t, createdQuestion.Created, data.Created)
	}
}

func TestQuestionGetByIDNotFound(t *testing.T) {
	cMock := getMock()
	cMock.On("GetQuestionByID", 0).Return(model.Question{}, ui.ErrNoResult)

	data, err := QuestionGET("0")
	if assert.NotNil(t, err) {
		cMock.AssertExpectations(t)

		assert.Nil(t, data)
		assert.Equal(t, ui.ErrNoResult, err)
	}
}

func TestQuestionGetByAuthorIDCorrectData(t *testing.T) {
	cMock := getMock()

	createdQuestions := make([]model.Question, 0)
	createdQuestions = append(createdQuestions, createdQuestion)
	createdQuestions = append(createdQuestions, createdQuestion)
	cMock.On("GetQuestionsByAuthorID", 1).Return(createdQuestions, nil)

	data, err := QuestionsGET("1")
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 2, len(data))
		for _, d := range data {
			assert.Equal(t, createdQuestion.ID, d.ID)
			assert.Equal(t, createdQuestion.Title, d.Title)
			assert.Equal(t, *createdQuestion.Content, *d.Content)
			assert.Equal(t, createdQuestion.AuthorID, d.AuthorID)
			assert.Equal(t, *createdQuestion.HasBest, *d.HasBest)
			assert.Equal(t, createdQuestion.Created, d.Created)
		}
	}
}

func TestQuestionGetByAuthorIDNotFound(t *testing.T) {
	cMock := getMock()

	createdQuestions := make([]model.Question, 0)
	cMock.On("GetQuestionsByAuthorID", 1).Return(createdQuestions, nil)

	data, err := QuestionsGET("1")
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)

		assert.Equal(t, 0, len(data))
		assert.Equal(t, make([]model.Question, 0), data)
	}
}

/*
********************************************************************
TESTS FOR UPDATE QUESTION ******************************************
********************************************************************
*/

func TestUpdateCorrectData(t *testing.T) {
	cMock := getMock()
	cMock.On("UpdateQuestion", updatedQuestion).Return(updatedQuestion, nil)

	body, _ := json.Marshal(updatedQuestion)
	response, err := UpdatePATCH(body)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)

		assert.Equal(t, *updatedQuestion.HasBest, *response.HasBest)
		assert.Equal(t, *updatedQuestion.Content, *response.Content)
	}
}

func TestUpdateNotFound(t *testing.T) {
	cMock := getMock()
	cMock.On("UpdateQuestion", updatedQuestion).Return(model.Question{}, ui.ErrNoDataToUpdate)

	body, _ := json.Marshal(updatedQuestion)
	data, err := UpdatePATCH(body)
	if assert.NotNil(t, err) {
		cMock.AssertExpectations(t)

		assert.Equal(t, ui.ErrNoDataToUpdate, err)
		assert.Nil(t, data)
	}
}

func TestUpdateMissedID(t *testing.T) {
	body := []byte("{\"has_best\": true, \"content\": \"My New Content\"}")

	response, err := UpdatePATCH(body)
	assert.Equal(t, ui.ErrFieldsRequired, err)
	assert.Equal(t, (*model.Question)(nil), response)
}

func TestUpdateBrokenBody(t *testing.T) {
	data, err := UpdatePATCH([]byte("{id: 1}"))

	if assert.NotNil(t, err) {
		assert.Nil(t, data)
	}
}

/*
********************************************************************
TESTS FOR REMOVE QUESTION ******************************************
********************************************************************
*/

func TestRemoveByIDCorrectData(t *testing.T) {
	cMock := getMock()
	cMock.On("DeleteQuestionByID", questionToRemoveID).Return(nil)

	body := []byte("{\"id\": 1}")
	err := RemoveDELETE(body)
	assert.Nil(t, err)
}

func TestRemoveByAuthorIDCorrectData(t *testing.T) {
	cMock := getMock()
	cMock.On("DeleteQuestionByAuthorID", questionToRemoveAuthorID).Return(nil)

	body := []byte("{\"author_id\": 1}")
	err := RemoveDELETE(body)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
	}
}

func TestRemoveByIDNotFound(t *testing.T) {
	cMock := getMock()
	cMock.On("DeleteQuestionByID", questionToRemoveID).Return(ui.ErrNoDataToDelete)

	body := []byte("{\"id\": 1}")
	err := RemoveDELETE(body)
	if assert.Equal(t, ui.ErrNoDataToDelete, err) {
		cMock.AssertExpectations(t)
	}
}

func TestRemoveMissedIDs(t *testing.T) {
	body := []byte("{\"has_best\": true}")
	err := RemoveDELETE(body)
	assert.Equal(t, ui.ErrFieldsRequired, err)
}
