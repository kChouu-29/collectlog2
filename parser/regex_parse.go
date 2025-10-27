// parser/regex_parse.go

package parser

import (
	"collectlogupdate/model"
	"regexp"
)

// Regex đã được sửa:
// Khớp IP - - [TIMESTAMP] "METHOD PATH PROTOCOL" STATUS SIZE "REFERER" "AGENT"
var LogRegex = regexp.MustCompile(`(?P<cline_ip>\S+) - - \[(?P<timestamp>[^\]]+)\] "(?P<method>\S+) (?P<path>\S+) \S+" (?P<status>\d{3}) (?P<size>\d+) "(?P<referer>[^"]*)" "(?P<agent>[^"]*)"`)

func ParserAccessLog(line string) (*model.Log, bool) {
	matches := LogRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil, false
	}
	// Chỉ mục (index) của matches[] được tính dựa trên Regex (1 = IP, 2 = Timestamp, v.v.)
	return &model.Log{
		ClineIP:   matches[1], // Index 1: IP
		Timestamp: matches[2], // Index 2: Timestamp
		Method:    matches[3],
		Path:      matches[4],
		Status:    matches[5],
		Size:      matches[6],
		// Index 7 là Referer (không có trong model.Log) nên ta bỏ qua
		Agent: matches[8], // Index 8: Agent
	}, true
}
