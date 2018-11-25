package entity

import "time"

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
