// handler/handler.go

package handler

import (
	"collectlogupdate/controller"
	"encoding/json" // <-- THÊM THƯ VIỆN NÀY
	"io"
	"log"
	"net/http"
)

// Struct để giải mã JSON từ Filebeat.
// Dòng log Nginx thực sự nằm trong trường "message"
type FilebeatEvent struct {
	Message string `json:"message"`
}

func ReceiveLog(w http.ResponseWriter, r *http.Request) {
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

	// Khai báo mảng để hứng các sự kiện (vì Filebeat thường gửi theo lô)
	var events []FilebeatEvent
	if err := json.Unmarshal(body, &events); err != nil {
		log.Println("Failed to unmarshal Filebeat JSON:", err)
		http.Error(w, "Failed to parse body as JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Received %d log events", len(events))

	// Lặp qua từng sự kiện và xử lý dòng log thô bên trong trường Message
	for _, event := range events {
		go controller.ProcessLog(event.Message)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Done"))
}