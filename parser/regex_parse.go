// parser/regex_parse.go
package parser

import (
	"collectlogupdate/model"
	"regexp"
	"strconv"
	"time"
)

var LogRegex = regexp.MustCompile(`(?P<cline_ip>\S+) - - \[(?P<timestamp>[^\]]+)\] "(?P<method>\S+) (?P<path>\S+) \S+" (?P<status>\d{3}) (?P<size>\d+) "(?P<referer>[^"]*)" "(?P<agent>[^"]*)"`)

func ParserAccessLog(line string) (*model.Log, bool) {
	matches := LogRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil, false
	}

	// Parse timestamp từ format nginx: 02/Jan/2006:15:04:05 -0700
	timestamp, err := time.Parse("02/Jan/2006:15:04:05 -0700", matches[2])
	if err != nil {
		// Fallback to current time if parse fails
		timestamp = time.Now()
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
