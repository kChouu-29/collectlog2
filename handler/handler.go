// handler/handler.go

package handler

import (
	"collectlogupdate/controller"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// Struct để giải mã JSON từ Filebeat (data part của Bulk API).
type FilebeatEvent struct {
	Message string `json:"message"`
	LogType string `json:"log_type"`
}

// Struct cho response của Elasticsearch-compatible API
type ElasticsearchResponse struct {
	Version struct {
		Number string `json:"number"`
	} `json:"version"`
	Tagline string `json:"tagline"`
}

func ReceiveLog(w http.ResponseWriter, r *http.Request) {
	// Xử lý GET hoặc HEAD request (Health Check)
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		w.Header().Set("Content-Type", "application/json")
		response := ElasticsearchResponse{
			Tagline: "You Know, for Search",
		}
		response.Version.Number = "8.15.0"
		json.NewEncoder(w).Encode(response)
		return
	}

	// Xử lý POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	var receivedEvents int

	// Loop để đọc các đối tượng JSON liên tiếp (NDJSON/Bulk API format)
	for {
		// 1. Đọc Metadata (ví dụ: { "index": {} })
		var meta map[string]interface{}
		if err := decoder.Decode(&meta); err == io.EOF {
			break // Kết thúc stream
		} else if err != nil {
			log.Println("Failed to decode Bulk API metadata:", err)
			break
		}

		// 2. Đọc Data (Log Event)
		var event FilebeatEvent
		if err := decoder.Decode(&event); err == io.EOF {
			log.Println("Unexpected EOF after reading metadata.")
			break
		} else if err != nil {
			log.Println("Failed to decode Bulk API event data:", err)
			break
		}

		receivedEvents++
		// Xử lý từng event
		go controller.ProcessLog(event.Message, event.LogType)
	}

	log.Printf("Received %d log events", receivedEvents)

	// Trả về JSON response giống Elasticsearch (200 OK)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"took":   1,
		"errors": false, // Báo không có lỗi để Filebeat không retry
		"items":  []interface{}{},
	})
}
