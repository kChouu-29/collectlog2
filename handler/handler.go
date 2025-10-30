//handler/handler.go

package handler

import (
	"collectlogupdate/controller"
	"collectlogupdate/model"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Struct để giải mã JSON từ Filebeat (data part của Bulk API).

func ReceiveLog(w http.ResponseWriter, r *http.Request) {
	// Xử lý GET hoặc HEAD request (Health Check)
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		w.Header().Set("Content-Type", "application/json")
		response := model.ElasticsearchResponse{
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

	// Tạo items array cho response
	items := []map[string]interface{}{}

	// Loop để đọc các đối tượng JSON liên tiếp (NDJSON/Bulk API format)
	for {
		// 1. Đọc Metadata
		var meta map[string]interface{}
		if err := decoder.Decode(&meta); err == io.EOF {
			break
		} else if err != nil {
			log.Println("Failed to decode Bulk API metadata:", err)
			break
		}

		// 2. Đọc Data (Log Event)
		var event model.FilebeatEvent
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

		// QUAN TRỌNG: Thêm item vào response để Filebeat biết đã xử lý thành công
		items = append(items, map[string]interface{}{
			"index": map[string]interface{}{
				"_index": "logs",
				"_id":    fmt.Sprintf("%d", receivedEvents),
				"status": 201,
				"result": "created",
			},
		})
	}

	log.Printf("Batch completed: Received %d log events", receivedEvents)

	// Trả về JSON response giống Elasticsearch
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"took":   1,
		"errors": false,
		"items":  items, // Phải trả về items để Filebeat biết đã xử lý
	})
}
