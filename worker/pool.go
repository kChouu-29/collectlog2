package worker

import (
	"collectlogupdate/controller"
	"collectlogupdate/model"
	"log"
)
const (
	NumWorkers = 10  
	QueueSize  = 1000 
)
var JobQueue chan model.FilebeatEvent

func StartWorkerPool() {

	JobQueue = make(chan model.FilebeatEvent, QueueSize) 

	for i := 1; i <= NumWorkers; i++ {
		go worker(i, JobQueue)
	}
	log.Printf("Started worker pool with %d workers and queue size %d", NumWorkers, QueueSize)
}

func worker(id int, jobs <-chan model.FilebeatEvent) {

	for job := range jobs {

		controller.ProcessLog(job.Message, job.LogType)
	}
}