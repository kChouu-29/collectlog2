package repository

import (
	"bytes"
	"collectlogupdate/model"
	"encoding/json"
	"log"
	"net/http"
)

func SaveToElastic(logData *model.Log) {
	data, _ := json.Marshal(logData)
	_, err := http.Post("http://elasticsearch:9200/logs/_doc", "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println("Error saving log to Elastic:", err)
	} else {
		log.Println("Log saved to Elastic successfully", logData.Path)
	}
}
