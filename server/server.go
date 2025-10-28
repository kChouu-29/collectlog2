package server

import (
	"collectlogupdate/handler"
	"log"
	"net/http"
)

func StartServer() {
	http.HandleFunc("/", handler.ReceiveLog)
	http.HandleFunc("/collect", handler.ReceiveLog)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
