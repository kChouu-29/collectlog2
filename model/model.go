package model

import "time"

type Log struct {
	Timestamp time.Time `json:"@timestamp"`
	ClientIP  string `json:"client.ip"`
	Method    string `json:"http.request.method"`
	Path      string `json:"http.request.path"`
	Status    int    `json:"http.response.status_code"`
	Size      int    `json:"http.response.body.bytes"`
	Agent     string `json:"agent"`
}
type FilebeatEvent struct {
	Message string `json:"message"`
	LogType string `json:"log_type"`
}


type ElasticsearchResponse struct {
	Version struct {
		Number string `json:"number"`
	} `json:"version"`
	Tagline string `json:"tagline"`
}