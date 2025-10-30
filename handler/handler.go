//handler/handler.go

package handler

import (
	"collectlogupdate/model"
	"collectlogupdate/worker"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)



func ReceiveLog(w http.ResponseWriter, r *http.Request) {
	
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		w.Header().Set("Content-Type", "application/json")
		response := model.ElasticsearchResponse{
			Tagline: "You Know, for Search",
		}
		response.Version.Number = "8.15.0"
		json.NewEncoder(w).Encode(response)
		return
	}

	
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	var receivedEvents int

	
	items := []map[string]interface{}{}

	
	for {
		
		var meta map[string]interface{}
		if err := decoder.Decode(&meta); err == io.EOF {
			break
		} else if err != nil {
			log.Println("Failed to decode Bulk API metadata:", err)
			break
		}

		
		var event model.FilebeatEvent
		if err := decoder.Decode(&event); err == io.EOF {
			log.Println("Unexpected EOF after reading metadata.")
			break
		} else if err != nil {
			log.Println("Failed to decode Bulk API event data:", err)
			break
		}

		receivedEvents++

		
		job := model.FilebeatEvent{
			Message: event.Message,
			LogType: event.LogType,
		}

		
		worker.JobQueue <- job

		// ======================

		// QUAN TRỌNG: Thêm item vào response để Filebeat biết đã xử lý thành công
		items = append(items, map[string]interface{}{
			"index": map[string]interface{}{
				"_index": "logs",
				"_id":    fmt.Sprintf("%d", receivedEvents),
				"status": 201, // 201 = Created
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
