package controller

import (
	"collectlogupdate/model"
	"collectlogupdate/parser"
	"collectlogupdate/repository"
	"log"
)

// <--- ĐÃ SỬA: Struct đại diện cho một công việc log
type LogJob struct {
	Line    string
	LogType string
}

// <--- ĐÃ SỬA: Kênh (channel) để chứa các công việc log đang chờ xử lý
var LogQueue chan LogJob

// <--- ĐÃ SỬA: Khởi động Worker Pool
func StartWorkerPool(numWorkers int) {
	// Khởi tạo channel với buffer size 1000 (cho phép 1000 log chờ trong hàng đợi)
	LogQueue = make(chan LogJob, 1000)
	log.Printf("Starting %d log processing workers...", numWorkers)

	// Khởi tạo số lượng worker cố định
	for i := 1; i <= numWorkers; i++ {
		go worker(i, LogQueue)
	}
}

// <--- ĐÃ SỬA: Hàm Worker chạy trong Goroutine
func worker(id int, jobs <-chan LogJob) {
	// Worker liên tục lấy job từ channel
	for job := range jobs {
		processLogCore(job.Line, job.LogType)
	}
}

// <--- ĐÃ SỬA: Logic xử lý log cũ được đổi tên thành processLogCore (logic chính)
func processLogCore(line string, logType string) {
	var parsed *model.Log
	var ok bool

	switch logType {
	case "access":
		parsed, ok = parser.ParserAccessLog(line)
	case "error":
		parsed, ok = parser.ParserErrorLog(line)
	default:
		log.Println("Unhandled log type:", logType)
		return
	}

	if !ok {
		log.Println("Failed to parse log line:", line)
		return
	}

	log.Println("Parsed log", parsed.Path)
	repository.SaveToElastic(parsed)
}

// <--- ĐÃ SỬA: Hàm ProcessLog mới chỉ có nhiệm vụ đưa công việc vào hàng đợi
func ProcessLog(line string, logType string) {
	// Log được đưa vào channel (không cần 'go' ở đây)
	// Channel sẽ giới hạn công việc và áp dụng backpressure nếu quá tải
	LogQueue <- LogJob{Line: line, LogType: logType}
}
