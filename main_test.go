package main

import (
	"testing"
	"time"
	"workerpool/pkg"
)

func testWorker(data string) interface{} {
    time.Sleep(100 * time.Millisecond)
    return data + "_processed"
}

func TestAddWorker(t *testing.T) {
    wp := pkg.CreateWorkerPool()

    err := wp.AddWorker("worker1", testWorker)
    if err != nil {
        t.Errorf("Error adding worker: %v", err)
    }

    err = wp.AddWorker("worker1", testWorker)
    if err == nil {
        t.Errorf("Expected error when adding worker with duplicate name, but got none")
    }
}

func TestRemoveWorker(t *testing.T) {
    wp := pkg.CreateWorkerPool()

    err := wp.AddWorker("worker1", testWorker)
    
    if err != nil {
        t.Errorf("Error adding worker: %v", err)
    }

    err = wp.RemoveWorker("worker1")
    if err != nil {
        t.Errorf("Error removing worker: %v", err)
    }

    err = wp.RemoveWorker("worker1")
    if err == nil {
        t.Errorf("Expected error when removing non-existent worker, but got none")
    }
}

func TestAddData(t *testing.T) {
    wp := pkg.CreateWorkerPool()
    wp.AddWorker("worker1", testWorker)
    wp.AddWorker("worker2", testWorker)

    wp.AddData("test1")
    wp.AddData("test2")

    time.Sleep(500 * time.Millisecond)

}

func TestBroadcastData(t *testing.T) {
    wp := pkg.CreateWorkerPool()
    wp.AddWorker("worker1", testWorker)
    wp.AddWorker("worker2", testWorker)

    wp.AddData("test1")
    wp.AddData("test2")

    time.Sleep(500 * time.Millisecond)

    wp.CloseWorkerPool()

}
