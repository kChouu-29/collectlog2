package main

import (
	"collectlogupdate/server"
	"collectlogupdate/worker"
)



func main() {

	worker.StartWorkerPool() 
	server.StartServer()
}