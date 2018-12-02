package main

import (
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/RSOI/question/controller"
	"github.com/RSOI/question/model"

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
	return args.Get(0).(model.Question), args.Error(0)
}
func (s *MockedQService) GetQuestionsByAuthorID(qAuthorID int) ([]model.Question, error) {
	args := s.Mock.Called(qAuthorID)
	return args.Get(0).([]model.Question), args.Error(0)
}
func (s *MockedQService) UpdateQuestion(q model.Question) (model.Question, error) {
	args := s.Mock.Called(q)
	return args.Get(0).(model.Question), args.Error(0)
}
func (s *MockedQService) GetUsageStatistic(host string) (model.ServiceStatus, error) {
	args := s.Mock.Called(host)
	return args.Get(0).(model.ServiceStatus), args.Error(0)
}
func (s *MockedQService) LogStat(request []byte, responseStatus int, responseError string) {

}

var (
	HOST                          = "http://localhost:8080"
	defaultQuestionContent        = "My Content"
	defaultQuestionHasBest        = false
	defaultQuestionCreatedTime, _ = time.Parse("2006-01-02T15:04:05", time.Now().String())
	defaultQuestion               = model.Question{
		AuthorID: 1,
		Title:    "My Title",
		Content:  &defaultQuestionContent,
	}
	createdQueston = model.Question{
		ID:       1,
		AuthorID: 1,
		Title:    "My Title",
		Content:  &defaultQuestionContent,
		HasBest:  &defaultQuestionHasBest,
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

func TestAskAddCorrectData(t *testing.T) {
	client, req, res, cMock := initServer()

	q, _ := json.Marshal(&defaultQuestion)

	req.SetRequestURI(HOST + "/ask")
	req.Header.SetMethod("PUT")
	req.SetBody(q)

	cMock.On("AddQuestion", defaultQuestion).Return(createdQueston, nil)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 201, res.Header.StatusCode())

		var response controller.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 201, response.Status)
		assert.Equal(t, "", response.Error)
		responseData := response.Data.(map[string]interface{})
		assert.Equal(t, createdQueston.ID, int(responseData["id"].(float64)))
		assert.Equal(t, createdQueston.Title, responseData["title"])
		assert.Equal(t, *createdQueston.Content, responseData["content"])
		assert.Equal(t, createdQueston.AuthorID, int(responseData["author_id"].(float64)))
		assert.Equal(t, *createdQueston.HasBest, responseData["has_best"])
		assert.Equal(t, createdQueston.Created.Format("2006-01-02T15:04:05")+"Z", responseData["created"])
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

		var response controller.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 400, response.Status)
		assert.Equal(t, "required fields are empty: title", response.Error)
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

		var response controller.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 400, response.Status)
		assert.Equal(t, "required fields are empty: title, content", response.Error)
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

		var response controller.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 500, response.Status)
		assert.NotEqual(t, "", response.Error)
		assert.Equal(t, nil, response.Data)
	}
}
