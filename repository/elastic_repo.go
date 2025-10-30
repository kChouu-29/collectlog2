package repository

import (
	"bytes"
	"collectlogupdate/model"
	"encoding/json"
	"io" 
	"log"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

func SaveToElastic(logData *model.Log) {
	data, _ := json.Marshal(logData)


	resp, err := http.Post("http://elasticsearch:9200/logs/_doc", "application/json", bytes.NewBuffer(data))

	if err != nil {
		
		log.Println("Error saving log to Elastic (Network Error):", err)
		return
	}

	
	defer resp.Body.Close()

	
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Println("Log saved to Elastic successfully", logData.Path)
	} else {
		
		errorBody, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			log.Printf("Error saving log to Elastic (Status %d) for path %s, but failed to read error body: %v", resp.StatusCode, logData.Path, readErr)
		} else {
			log.Printf("Error saving log to Elastic (Status %d) for path %s. Response body: %s", resp.StatusCode, logData.Path, string(errorBody))
		}
	}
}
