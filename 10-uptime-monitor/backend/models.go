package main

type Service struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	URL             string `json:"url"`
	IntervalSeconds int    `json:"interval_seconds"`
	CreatedAt       string `json:"created_at"`
}

type ServiceRequest struct {
	Name            string `json:"name"`
	URL             string `json:"url"`
	IntervalSeconds int    `json:"interval_seconds"`
}

type Check struct {
	ID             int    `json:"id"`
	ServiceID      int    `json:"service_id"`
	Status         string `json:"status"`
	ResponseTimeMs int    `json:"response_time_ms"`
	CheckedAt      string `json:"checked_at"`
}
