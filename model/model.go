package model

type Log struct {
	Timestamp string `json:"@timestamp"`
	ClineIP string `json:"cline.ip"`
	Method string `json:"http.request.method"`
	Path string `json:"http.request.path"`
	Status string `json:"http.response.status_code"`
	Size string `json:"http.response.body.bytes"`
	Agent string `json:"agent"`
}