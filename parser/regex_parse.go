// parser/regex_parse.go
package parser

import (
	"collectlogupdate/model"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// <--- ĐÃ SỬA: Đổi tên LogRegex thành AccessLogRegex
var AccessLogRegex = regexp.MustCompile(`(?P<cline_ip>\S+) - - \[(?P<timestamp>[^\]]+)\] "(?P<method>\S+) (?P<path>\S+) \S+" (?P<status>\d{3}) (?P<size>\d+) "(?P<referer>[^"]*)" "(?P<agent>[^"]*)"`)

// <--- ĐÃ SỬA: Regex mới cho Nginx Error Log
var ErrorLogRegex = regexp.MustCompile(`^(?P<timestamp>\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}) \[(?P<level>\w+)\] .*?: (?P<message>.*)`)

func ParserAccessLog(line string) (*model.Log, bool) {
	matches := AccessLogRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil, false
	}

	// Parse timestamp từ format nginx: 02/Jan/2006:15:04:05 -0700
	timestamp, err := time.Parse("02/Jan/2006:15:04:05 -0700", matches[2])
	if err != nil {
		// <--- ĐÃ SỬA: Trả về false thay vì fallback về time.Now()
		return nil, false
	}

	status, _ := strconv.Atoi(matches[5])
	size, _ := strconv.Atoi(matches[6])

	return &model.Log{
		Timestamp: timestamp,
		ClientIP:  matches[1],
		Method:    matches[3],
		Path:      matches[4],
		Status:    status,
		Size:      size,
		Agent:     matches[8],
	}, true
}

// <--- ĐÃ SỬA: Hàm parser mới cho Log Lỗi
func ParserErrorLog(line string) (*model.Log, bool) {
	matches := ErrorLogRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil, false
	}

	// Parse timestamp từ format: 2006/01/02 15:04:05
	timestamp, err := time.Parse("2006/01/02 15:04:05", matches[1])
	if err != nil {
		return nil, false
	}
	
	level := matches[2]
	message := matches[3]

	// Map thông tin log lỗi vào model.Log hiện có.
	// Sử dụng trường Path để lưu thông tin loại log và message.
	return &model.Log{
		Timestamp: timestamp,
		ClientIP:  "0.0.0.0", // Giả định IP 0.0.0.0 vì log lỗi không luôn có client IP.
		Method:    "ERROR_LOG",
		Path:      "[" + strings.ToUpper(level) + "]: " + message,
		Status:    500, // Gán status 500 cho các lỗi để dễ dàng lọc.
		Size:      0,
		Agent:     "Nginx/ErrorLog",
	}, true
}