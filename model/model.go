package model

import (
	"time"
)

type FetchRequest struct {
	StartDate string  `json:"startDate"  validate:"required"`
	EndDate   string  `json:"endDate"  validate:"required"`
	MinCount  float64 `json:"minCount"  validate:"required"`
	MaxCount  float64 `json:"maxCount"  validate:"required"`
}

type FetchResponse struct {
	Code    int       `json:"code"`
	Msg     string    `json:"msg"`
	Records []Records `json:"records"`
}
type Records struct {
	Key        string    `json:"key"`
	CreatedAt  time.Time `json:"createdAt"`
	TotalCount int       `json:"totalCount"`
}

type InMemoryDataIO struct {
	Key   string `json:"key"  validate:"required"`
	Value string `json:"value" validate:"required"`
}
type ApiError struct {
	Message string `json:"message"`
}
