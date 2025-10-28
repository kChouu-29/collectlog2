// handler/handler.go

package handler

import (
	"collectlogupdate/controller"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// Struct để giải mã JSON từ Filebeat.
type FilebeatEvent struct {
	Message string `json:"message"`
}

// Struct cho response của Elasticsearch-compatible API
type ElasticsearchResponse struct {
	Version struct {
		Number string `json:"number"`
	} `json:"version"`
	Tagline string `json:"tagline"`
}

func ReceiveLog(w http.ResponseWriter, r *http.Request) {
	// Nếu là GET request, trả về thông tin health check (giả lập Elasticsearch)
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		response := ElasticsearchResponse{
			Tagline: "You Know, for Search",
		}
		response.Version.Number = "8.15.0"
		json.NewEncoder(w).Encode(response)
		return
	}

	// Xử lý POST request như cũ
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Khai báo mảng để hứng các sự kiện
	var events []FilebeatEvent
	if err := json.Unmarshal(body, &events); err != nil {
		log.Println("Failed to unmarshal Filebeat JSON:", err)
		http.Error(w, "Failed to parse body as JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Received %d log events", len(events))

	// Xử lý từng event
	for _, event := range events {
		go controller.ProcessLog(event.Message)
	}

	// Trả về JSON response giống Elasticsearch
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"took":   1,
		"errors": false,
		"items":  []interface{}{},
	})
}
