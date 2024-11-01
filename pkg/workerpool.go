package main

import (
	"log"
	"sync"
	"time"
)

type WorkerPool struct {
	workerCount int
	DataChan    chan string
	wg          sync.WaitGroup
	stopChan	chan string
}

func (wp *WorkerPool) AddWorker(name string, worker func(string)) {
	wp.wg.Add(1)

	go func(name string) {
		defer wp.wg.Done()

		fmt.Printf()
		
		for command := range wg.DataChan {
			select {

			case stopedFuncName := <-wp.stopChan:
				if stopedFuncName == name {
					return
				}

			default:
				start := time.Now()
				worker(name)
				log.Printf("")
			}
		}
	}
}

func (wp * WorkerPool) RemoveWorker(name string) {

}

func main() {

}