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
	Workers 	[]string
	chanels     map[string]chan string
	mu 			*sync.Mutex
}

func (wp *WorkerPool) AddWorker(name string, worker func(string) interface{}) error {

	if slices.Contains(wp.Workers, name) {
		err := fmt.Errorf("worker with name %s already exist, please change name and try again", name)
		log.Printf("Error: %s", err)
		return err
	}
	
	defer wp.wg.Done()

	wp.mu.Lock()

	workerChan := make(chan string)
    wp.chanels[name] = workerChan

	wp.wg.Add(1)
	go func(name string, worker func(string) interface{}) {
		
		for command := range workerChan {

			start := time.Now()
			result := worker(command)
			log.Printf("Worker '%s' processed the string '%s' with result '%s' in %s\n", name, command, result, time.Since(start))
			wp.wg.Done()
		}
	}(name, worker)

	wp.Workers = append(wp.Workers, name)
	log.Printf("Worker %s started\n", name)
	wp.mu.Unlock()

	return nil
}

func (wp *WorkerPool) RemoveWorker(name string) error {

	if !slices.Contains(wp.Workers, name) {
		err := fmt.Errorf("worker with name %s not exist, please change name and try again", name)
		log.Printf("Error: %s", err)
		return err
	}

	defer wp.wg.Done()
	wp.wg.Add(1)

	wp.mu.Lock()

	close(wp.chanels[name])
	delete(wp.chanels, name)

	for idx, i := range wp.Workers {
		if i == name {
			wp.Workers[idx] = wp.Workers[len(wp.Workers)-1]
			wp.Workers = wp.Workers[:len(wp.Workers)-1]
			break
		}
	}

	wp.mu.Unlock()
	log.Printf("Worker %s stopped\n", name)
	
	return nil
}

func (wp *WorkerPool) BroadcastData() {

	for data := range wp.DataChan {
		wp.wg.Add(len(wp.chanels))

		go func () {
			for _, v := range wp.chanels {
				v <-data
		}
		}()
	}
}

func (wp *WorkerPool) CloseWorkerPool() {
	time.Sleep(time.Millisecond)
	wp.wg.Wait()

	close(wp.DataChan)
	for _, v := range wp.chanels {
		close(v)
	}

	log.Println("Worker-pool closed")
}

func (wp *WorkerPool) AddData(name string) {
	wp.DataChan<- name
}

func CreateWorkerPool() *WorkerPool {
	workerPool := WorkerPool{
		DataChan: make(chan string),
		wg: &sync.WaitGroup{},
		Workers: []string{},
		chanels: make(map[string]chan string),
		mu: &sync.Mutex{},
	}
	go workerPool.BroadcastData()

	log.Println("New worker-pool created")
	return &workerPool
}