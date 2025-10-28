package repository

import (
	"bytes"
	"collectlogupdate/model"
	"encoding/json"
	"io" // <--- ĐÃ SỬA: Bổ sung import io để đọc body response
	"log"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 5 * time.Second, // Đặt giới hạn 5 giây
}

func SaveToElastic(logData *model.Log) {
	data, _ := json.Marshal(logData)

	// <--- ĐÃ SỬA: Lấy response và error
	resp, err := http.Post("http://elasticsearch:9200/logs/_doc", "application/json", bytes.NewBuffer(data))

	if err != nil {
		// 1. Lỗi Cấp độ Mạng (ví dụ: không kết nối được đến Elasticsearch)
		log.Println("Error saving log to Elastic (Network Error):", err)
		return
	}

	// Đảm bảo đóng Body của response để tránh rò rỉ tài nguyên
	defer resp.Body.Close()

	// 2. Lỗi Cấp độ Ứng dụng (Kiểm tra Status Code)
	// Mã trạng thái 2xx là thành công
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Println("Log saved to Elastic successfully", logData.Path)
	} else {
		// Log lỗi chi tiết nếu status code không phải 2xx
		errorBody, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			log.Printf("Error saving log to Elastic (Status %d) for path %s, but failed to read error body: %v", resp.StatusCode, logData.Path, readErr)
		} else {
			log.Printf("Error saving log to Elastic (Status %d) for path %s. Response body: %s", resp.StatusCode, logData.Path, string(errorBody))
		}
	}
}
