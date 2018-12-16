package controller

import (
	"github.com/RSOI/question/model"
)

// IndexGET returns usage statistic
func IndexGET(host []byte) (*model.ServiceStatus, error) {
	data, err := QuestionModel.GetUsageStatistic(string(host))
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// LogStat stores service usage
func LogStat(path []byte, status int, err string) {
	QuestionModel.LogStat(path, status, err)
}
