package handler

import (
	"collectlogupdate/controller"
	"io"
	"log"
	"net/http"
)

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

	log.Println("Received Log")

	go controller.ProcessLog(string(body))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Done"))
}
