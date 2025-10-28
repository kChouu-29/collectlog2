package controller

import (
	"collectlogupdate/model"
	"collectlogupdate/parser"
	"collectlogupdate/repository"
	"log"
)

// Hàm ProcessLog mới chứa logic xử lý log trực tiếp.
// Hàm này được gọi trong một Goroutine từ Handler.
func ProcessLog(line string, logType string) {
	var parsed *model.Log
	var ok bool

	// Lựa chọn parser dựa trên loại log
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
	// Gọi hàm lưu trữ
	repository.SaveToElastic(parsed)
}