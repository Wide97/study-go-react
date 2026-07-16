package main

import "time"

type Service struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	URL             string    `json:"url"`
	IntervalSeconds int       `json:"interval_seconds"`
	CreatedAt       time.Time `json:"created_at"`
}
