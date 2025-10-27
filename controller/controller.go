package controller

import (
	"collectlogupdate/parser"
	"collectlogupdate/repository"
	"log"
)

func ProcessLog(line string) {
	parsed, ok := parser.ParserAccessLog(line)
	if !ok {
		return
	}
	log.Println("Parsed log", parsed.Path)
	repository.SaveToElastic(parsed)
}
