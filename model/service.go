package model

import (
	"time"

	"github.com/RSOI/question/database"
	"github.com/jackc/pgx"
)

// RequestInfo interface. Provides usage data.
type RequestInfo struct {
	Request     string    `json:"request"`
	RequestTime time.Time `json:"request_time"`
}

// ServiceStatus interface. Provides usage data.
type ServiceStatus struct {
	Address       string    `json:"address"`
	RequestsCount int       `json:"requests_count"`
	LastUsage     time.Time `json:"last_usage"`
}

// GetUsageStatistic provides access to logs
func GetUsageStatistic() ServiceStatus {
	row := database.DB.QueryRow(`
		SELECT cnt.*, last_usage.* FROM
			(SELECT count(*) FROM question.services) AS cnt,
			(SELECT RequestTime FROM question.services LIMIT 1) AS last_usage
	`)
	var ServiceResponse ServiceStatus
	err := row.Scan(&ServiceResponse.RequestsCount)
	switch err {
	case pgx.ErrNoRows:
		ServiceResponse.RequestsCount = 0
		ServiceResponse.LastUsage = time.Time{}
	case nil:
	default:
		panic(err)
	}

	return ServiceResponse
}
