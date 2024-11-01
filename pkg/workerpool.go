package pkg

import (
	"fmt"
	"log"
	"slices"
	"sync"
	"time"
)

type WorkerPool struct {
	DataChan    chan string
	wg          *sync.WaitGroup
	stopChan	chan string
	Workers 	[]string
}

func (wp *WorkerPool) AddWorker(name string, worker func(string)) error {

	if slices.Contains(wp.Workers, name) {
		return fmt.Errorf("worker with name %s already exist, please change name and try again", name)
	}
	
	wp.wg.Add(1)

	go func(name string, worker func(string)) {
		defer wp.wg.Done()
		
		for command := range wp.DataChan {
			select {

			case stopedFuncName := <-wp.stopChan:
				if stopedFuncName == name {
					return
				}

			default:
				start := time.Now()
				worker(command)
				log.Printf("Worker %s processed the string %s in %s\n", name, command, time.Since(start))
			}
		}
	}(name, worker)

	wp.Workers = append(wp.Workers, name)
	log.Printf("Worker %s started\n", name)

	return nil
}

func (wp *WorkerPool) RemoveWorker(name string) error {

	if !slices.Contains(wp.Workers, name) {
		return fmt.Errorf("worker with name %s not exist, please change name and try again", name)
	}

	wp.stopChan<- name
	log.Printf("Worker %s stopped\n", name)
	
	return nil
}

func (wp *WorkerPool) CloseWorkerPool() {
	close(wp.DataChan)
	wp.wg.Wait()
	log.Println("Worker-pool closed")
}

func CreateWorkerPool() *WorkerPool {
	workerPool := WorkerPool{
		DataChan: make(chan string),
		wg: &sync.WaitGroup{},
		stopChan: make(chan string),
		Workers: []string{},
	}

	log.Println("New worker-pool created")
	return &workerPool
}