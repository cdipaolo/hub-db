package model

import "time"

// HealthCheckResponse returns status info to check
// if the crawler is working.
type HealthCheckResponse struct {
	Status       string        `json:"status"`
	SaveTo       string        `json:"saving_json_to,omitempty"`
	Uptime       time.Duration `json:"uptime,omitempty"`
	PagesCrawled uint64        `json:"pages_crawled"`
	Errors       []error       `json:"errors"`
}
