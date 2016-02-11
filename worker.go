package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/if1live/haru/gallery"
)

// http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html

type WorkRequest struct {
	Service string
	Id      string
}

func (w WorkRequest) Run() {
	g := gallery.New(w.Service)
	g.PrefetchCover(w.Id)
	g.PrefetchImage(w.Id)
}

var WorkQueue = make(chan WorkRequest, 100)

func Collector(w http.ResponseWriter, r *http.Request, service string, id string) {
	work := WorkRequest{Service: service, Id: id}

	WorkQueue <- work
	log.Printf("Work request queued %q", work)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{status:"ok"}`)
	return
}

type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
}

func NewWorker(id int, workerQueue chan chan WorkRequest) Worker {
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
	}
	return worker
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				log.Printf("worker%d: Received work request", w.ID)
				work.Run()
				log.Printf("worker%d: work complete", w.ID)
			case <-w.QuitChan:
				log.Printf("worker%d stopping", w.ID)
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

var WorkerQueue chan chan WorkRequest

func StartDispatcher(nworkers int) {
	WorkerQueue = make(chan chan WorkRequest, nworkers)

	for i := 0; i < nworkers; i++ {
		log.Println("Starting worker", i+1)
		worker := NewWorker(i+1, WorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-WorkQueue:
				log.Println("Received work request")
				go func() {
					worker := <-WorkerQueue
					log.Println("Dispatching work request")
					worker <- work
				}()
			}
		}
	}()
}
