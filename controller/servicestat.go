package controller

import (
	"fmt"

	"github.com/RSOI/question/model"
	"github.com/RSOI/question/utils"
)

// IndexGET returns usage statistic
func IndexGET(host []byte) (*model.ServiceStatus, error) {
	data, err := QuestionModel.GetUsageStatistic(string(host))
	if err != nil {
		utils.LOG(fmt.Sprintf("Data error: %s", err.Error()))
		return nil, err
	}

	utils.LOG("Successfull accessing usage statistic")
	return &data, nil
}

// LogStat stores service usage
func LogStat(path []byte, status int, err string) {
	utils.LOG("Storing usage stat...")
	QuestionModel.LogStat(path, status, err)
}
