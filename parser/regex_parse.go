package parser

import (
	"collectlogupdate/model"
	"regexp"
)

var LogRegex = regexp.MustCompile(`(?P<timestamp>\S+) - - \[(?P<cline_ip>[^\]]+)\] "(?P<method>\S+) (?P<path>\S+) \S+" (?P<status>\d{3}) (?P<size>\d+) "(?P<agent>[^"]+)"`)

func ParserAccessLog(line string) (*model.Log, bool) {
	matches := LogRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil, false
	}
	return &model.Log{
		Timestamp: matches[1],
		ClineIP:   matches[2],
		Method:    matches[3],
		Path:      matches[4],
		Status:    matches[5],
		Size:      matches[6],
		Agent:     matches[7],
	}, true
}
