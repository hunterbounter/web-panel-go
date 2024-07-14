package worker

import (
	"log"
	"sync"
	"time"
)

/*
Author: hunterbounter
This package can be used to start and control two different worker functions (workerA and workererB) that run at specific intervals.
Each worker runs for a given control time and stops running when a signal is received via the stop channel (StopChan). //canceled
The state of the workers is checked periodically in the main loop and if a worker has stopped, it is restarted. This package is suitable for situations where you want to have multiple worker threads running at the same time, each with different functions.
*/

// Worker struct
type Worker struct {
	Name          string
	CheckInterval time.Duration
	WorkerFunc    func(Worker)
}

/*
This Function Manages the User's Docker Processes.
If no Feed Data has been sent to the relevant Docker for 5 minutes, the Docker Process is terminated.
*/
func dockerManager(w Worker) {
	for {
		log.Println("Docker Manager Worker is Running")
		workerDockerManager()
		time.Sleep(w.CheckInterval)
	}

}

func Init() {
	log.Println("init worker")

	var wg sync.WaitGroup

	// Worker
	workers := []Worker{
		{Name: "Docker Manager", CheckInterval: 3 * time.Second, WorkerFunc: dockerManager},
	}

	for i := range workers {

		go workers[i].WorkerFunc(workers[i])
	}

	// Wait for workers to finish before the program ends
	wg.Wait()

}
