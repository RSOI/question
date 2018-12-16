package main

import (
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/RSOI/question/controller"
	"github.com/RSOI/question/model"
	"github.com/RSOI/question/ui"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedQService struct {
	mock.Mock
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
	HOST                          = "http://localhost"
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
)

func initServer() (*fasthttp.Client, *fasthttp.Request, *fasthttp.Response, *MockedQService) {
	listener := fasthttputil.NewInmemoryListener()
	server := &fasthttp.Server{
		Handler: initRoutes().Handler,
	}
	go server.Serve(listener)

	client := &fasthttp.Client{
		Dial: func(addr string) (net.Conn, error) {
			return listener.Dial()
		},
	}

	controller.QuestionModel = &MockedQService{}
	cMock := controller.QuestionModel.(*MockedQService)
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	return client, req, res, cMock
}

/*
********************************************************************
TESTS FOR ASK ******************************************************
********************************************************************
*/

func TestAskAddCorrectData(t *testing.T) {
	client, req, res, cMock := initServer()

	q, _ := json.Marshal(&defaultQuestion)

	req.SetRequestURI(HOST + "/ask")
	req.Header.SetMethod("PUT")
	req.SetBody(q)

	cMock.On("AddQuestion", defaultQuestion).Return(createdQuestion, nil)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 201, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 201, response.Status)
		assert.Equal(t, "", response.Error)
		responseData := response.Data.(map[string]interface{})
		assert.Equal(t, createdQuestion.ID, int(responseData["id"].(float64)))
		assert.Equal(t, createdQuestion.Title, responseData["title"])
		assert.Equal(t, *createdQuestion.Content, responseData["content"])
		assert.Equal(t, createdQuestion.AuthorID, int(responseData["author_id"].(float64)))
		assert.Equal(t, *createdQuestion.HasBest, responseData["has_best"])
		assert.Equal(t, createdQuestion.Created.Format("2006-01-02T15:04:05")+"Z", responseData["created"])
	}
}

func TestAskAddMissedOneField(t *testing.T) {
	client, req, res, _ := initServer()

	req.SetRequestURI(HOST + "/ask")
	req.Header.SetMethod("PUT")
	req.SetBodyString("{\"author_id\": 1, \"content\": \"My Content\"}")

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		assert.Equal(t, 400, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 400, response.Status)
		assert.Equal(t, ui.ErrFieldsRequired.Error(), response.Error)
		assert.Equal(t, nil, response.Data)
	}
}

func TestAskAddMissedManyField(t *testing.T) {
	client, req, res, _ := initServer()

	req.SetRequestURI(HOST + "/ask")
	req.Header.SetMethod("PUT")
	req.SetBodyString("{\"author_id\": 1}")

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		assert.Equal(t, 400, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 400, response.Status)
		assert.Equal(t, ui.ErrFieldsRequired.Error(), response.Error)
		assert.Equal(t, nil, response.Data)
	}
}

func TestAskAddBrokenBody(t *testing.T) {
	client, req, res, _ := initServer()

	req.SetRequestURI(HOST + "/ask")
	req.Header.SetMethod("PUT")
	req.SetBodyString("{author_id: 1}")

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		assert.Equal(t, 500, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 500, response.Status)
		assert.NotEqual(t, "", response.Error)
		assert.Equal(t, nil, response.Data)
	}
}

/*
********************************************************************
TESTS FOR QUESTION ID **********************************************
********************************************************************
*/

func TestQuestionGetByIDCorrectData(t *testing.T) {
	client, req, res, cMock := initServer()

	req.SetRequestURI(HOST + "/question/id1")
	req.Header.SetMethod("GET")

	cMock.On("GetQuestionByID", 1).Return(createdQuestion, nil)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 200, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 200, response.Status)
		assert.Equal(t, "", response.Error)
		responseData := response.Data.(map[string]interface{})
		assert.Equal(t, createdQuestion.ID, int(responseData["id"].(float64)))
		assert.Equal(t, createdQuestion.Title, responseData["title"])
		assert.Equal(t, *createdQuestion.Content, responseData["content"])
		assert.Equal(t, createdQuestion.AuthorID, int(responseData["author_id"].(float64)))
		assert.Equal(t, *createdQuestion.HasBest, responseData["has_best"])
		assert.Equal(t, createdQuestion.Created.Format("2006-01-02T15:04:05")+"Z", responseData["created"])
	}
}

func TestQuestionGetByIDNotFound(t *testing.T) {
	client, req, res, cMock := initServer()

	req.SetRequestURI(HOST + "/question/id0")
	req.Header.SetMethod("GET")

	cMock.On("GetQuestionByID", 0).Return(model.Question{}, ui.ErrNoResult)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 404, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 404, response.Status)
		assert.Equal(t, ui.ErrNoResult.Error(), response.Error)
		assert.Equal(t, nil, response.Data)
	}
}

/*
********************************************************************
TESTS FOR QUESTION AUTHORID ****************************************
********************************************************************
*/

func TestQuestionGetByAuthorIDCorrectData(t *testing.T) {
	client, req, res, cMock := initServer()

	req.SetRequestURI(HOST + "/questions/author1")
	req.Header.SetMethod("GET")

	createdQuestions := make([]model.Question, 0)
	createdQuestions = append(createdQuestions, createdQuestion)
	createdQuestions = append(createdQuestions, createdQuestion)
	cMock.On("GetQuestionsByAuthorID", 1).Return(createdQuestions, nil)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 200, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 200, response.Status)
		assert.Equal(t, "", response.Error)
		responseData := response.Data.([]interface{})
		assert.Equal(t, 2, len(responseData))
		for _, d := range responseData {
			data := d.(map[string]interface{})
			assert.Equal(t, createdQuestion.ID, int(data["id"].(float64)))
			assert.Equal(t, createdQuestion.Title, data["title"])
			assert.Equal(t, *createdQuestion.Content, data["content"])
			assert.Equal(t, createdQuestion.AuthorID, int(data["author_id"].(float64)))
			assert.Equal(t, *createdQuestion.HasBest, data["has_best"])
			assert.Equal(t, createdQuestion.Created.Format("2006-01-02T15:04:05")+"Z", data["created"])
		}
	}
}

func TestQuestionGetByAuthorIDNotFound(t *testing.T) {
	client, req, res, cMock := initServer()

	req.SetRequestURI(HOST + "/questions/author1")
	req.Header.SetMethod("GET")

	createdQuestions := make([]model.Question, 0)
	cMock.On("GetQuestionsByAuthorID", 1).Return(createdQuestions, nil)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 200, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 200, response.Status)
		assert.Equal(t, "", response.Error)
		responseData := response.Data.([]interface{})
		assert.Equal(t, 0, len(responseData))
		assert.Equal(t, make([]interface{}, 0), response.Data)
	}
}

/*
********************************************************************
TESTS FOR UPDATE QUESTION ******************************************
********************************************************************
*/

func TestUpdateCorrectData(t *testing.T) {
	client, req, res, cMock := initServer()

	source, _ := json.Marshal(updatedQuestion)

	req.SetRequestURI(HOST + "/update")
	req.Header.SetMethod("PATCH")
	req.SetBody(source)

	cMock.On("UpdateQuestion", updatedQuestion).Return(updatedQuestion, nil)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 200, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 200, response.Status)
		assert.Equal(t, "", response.Error)

		responseData := response.Data.(map[string]interface{})
		assert.Equal(t, *updatedQuestion.HasBest, responseData["has_best"])
		assert.Equal(t, *updatedQuestion.Content, responseData["content"])
	}
}

func TestUpdateNotFound(t *testing.T) {
	client, req, res, cMock := initServer()

	source, _ := json.Marshal(updatedQuestion)

	req.SetRequestURI(HOST + "/update")
	req.Header.SetMethod("PATCH")
	req.SetBody(source)

	cMock.On("UpdateQuestion", updatedQuestion).Return(model.Question{}, ui.ErrNoDataToUpdate)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 404, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 404, response.Status)
		assert.Equal(t, ui.ErrNoDataToUpdate.Error(), response.Error)
		assert.Equal(t, nil, response.Data)
	}
}

func TestUpdateMissedID(t *testing.T) {
	client, req, res, _ := initServer()

	req.SetRequestURI(HOST + "/update")
	req.Header.SetMethod("PATCH")
	req.SetBodyString("{\"has_best\": true, \"content\": \"My New Content\"}")

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		assert.Equal(t, 400, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 400, response.Status)
		assert.Equal(t, "missed required field(s)", response.Error)
		assert.Equal(t, nil, response.Data)
	}
}

func TestUpdateBrokenBody(t *testing.T) {
	client, req, res, _ := initServer()

	req.SetRequestURI(HOST + "/update")
	req.Header.SetMethod("PATCH")
	req.SetBodyString("{id: 1}")

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		assert.Equal(t, 500, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 500, response.Status)
		assert.NotEqual(t, "", response.Error)
		assert.Equal(t, nil, response.Data)
	}
}

/*
********************************************************************
TESTS FOR REMOVE QUESTION ******************************************
********************************************************************
*/

func TestRemoveByIDCorrectData(t *testing.T) {
	client, req, res, cMock := initServer()

	questionToRemove := model.Question{
		ID: 1,
	}
	questionToRemoveJSON, _ := json.Marshal(&questionToRemove)

	req.SetRequestURI(HOST + "/delete")
	req.Header.SetMethod("DELETE")
	req.SetBody(questionToRemoveJSON)

	cMock.On("DeleteQuestionByID", questionToRemove).Return(nil)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 200, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 200, response.Status)
		assert.Equal(t, "", response.Error)
		assert.Equal(t, nil, response.Data)
	}
}

func TestRemoveByAuthorIDCorrectData(t *testing.T) {
	client, req, res, cMock := initServer()

	questionToRemove := model.Question{
		AuthorID: 1,
	}
	questionToRemoveJSON, _ := json.Marshal(&questionToRemove)

	req.SetRequestURI(HOST + "/delete")
	req.Header.SetMethod("DELETE")
	req.SetBody(questionToRemoveJSON)

	cMock.On("DeleteQuestionByAuthorID", questionToRemove).Return(nil)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 200, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 200, response.Status)
		assert.Equal(t, "", response.Error)
		assert.Equal(t, nil, response.Data)
	}
}

func TestRemoveByIDNotFound(t *testing.T) {
	client, req, res, cMock := initServer()

	questionToRemove := model.Question{
		ID: 1,
	}
	questionToRemoveJSON, _ := json.Marshal(&questionToRemove)

	req.SetRequestURI(HOST + "/delete")
	req.Header.SetMethod("DELETE")
	req.SetBody(questionToRemoveJSON)

	cMock.On("DeleteQuestionByID", questionToRemove).Return(ui.ErrNoDataToDelete)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 404, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 404, response.Status)
		assert.Equal(t, ui.ErrNoDataToDelete.Error(), response.Error)
		assert.Equal(t, nil, response.Data)
	}
}

func TestRemoveMissedIDs(t *testing.T) {
	client, req, res, cMock := initServer()

	questionToRemove, _ := json.Marshal(&model.Question{})

	req.SetRequestURI(HOST + "/delete")
	req.Header.SetMethod("DELETE")
	req.SetBody(questionToRemove)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 400, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 400, response.Status)
		assert.Equal(t, ui.ErrFieldsRequired.Error(), response.Error)
		assert.Equal(t, nil, response.Data)
	}
}
