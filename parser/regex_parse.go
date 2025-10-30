// parser/regex_parse.go
package parser

import (
	"collectlogupdate/model"
	"regexp"
	"strconv"
	"strings"
	"time"
)


var AccessLogRegex = regexp.MustCompile(`(?P<cline_ip>\S+) - - \[(?P<timestamp>[^\]]+)\] "(?P<method>\S+) (?P<path>\S+) \S+" (?P<status>\d{3}) (?P<size>\d+) "(?P<referer>[^"]*)" "(?P<agent>[^"]*)"`)


var ErrorLogRegex = regexp.MustCompile(`^(?P<timestamp>\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}) \[(?P<level>\w+)\] .*?: (?P<message>.*)`)

func ParserAccessLog(line string) (*model.Log, bool) {
	matches := AccessLogRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil, false
	}

	timestamp, err := time.Parse("02/Jan/2006:15:04:05 -0700", matches[2])
	if err != nil {
	
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


func ParserErrorLog(line string) (*model.Log, bool) {
	matches := ErrorLogRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil, false
	}

	
	timestamp, err := time.Parse("2006/01/02 15:04:05", matches[1])
	if err != nil {
		return nil, false
	}
	
	level := matches[2]
	message := matches[3]

	
	return &model.Log{
		Timestamp: timestamp,
		ClientIP:  "0.0.0.0",
		Method:    "ERROR_LOG",
		Path:      "[" + strings.ToUpper(level) + "]: " + message,
		Status:    500, 
		Size:      0,
		Agent:     "Nginx/ErrorLog",
	}, true
}