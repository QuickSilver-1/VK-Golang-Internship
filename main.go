package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"workerpool/pkg"
)


func lenner(str string) interface{} {
	return strconv.Itoa(len(str))
}

func doubler(str string) interface{} {
	return str + str
}

func upper(str string) interface{} {
	return strings.ToUpper(str)
}

func lower(str string) interface{} {
	return strings.ToLower(str)
}

func longTime(str string) interface{} {
    time.Sleep(100 * time.Millisecond)
	return str
}

func main() {
	file, err := os.OpenFile("log/worker.log", os.O_APPEND, 0666)

	log.SetFlags(log.Ldate | log.Ltime)
	if err != nil {
		file, _ := os.OpenFile("log/file.log", os.O_APPEND, 0666)
		log.SetOutput(file)
		log.Fatal("Failed to open log file: ", err)
	}

	log.SetOutput(file)

	wp := pkg.CreateWorkerPool()
	
	wp.AddWorker("len", lenner)
	wp.AddWorker("doubler", doubler)
	wp.AddWorker("doubler", upper)
	wp.AddWorker("upper", upper)
	wp.AddWorker("lower", lower)
	wp.AddWorker("time", longTime)
	wp.RemoveWorker("1")

	for i := 0; i < 5; i++ {
		go wp.AddData(fmt.Sprintf("TeSt %d", i))
	}

	wp.CloseWorkerPool()
}