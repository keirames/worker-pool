package main

import (
	"context"
	"fmt"
	"sync"
)

type Job struct {
	content string
	// execFn  func(ctx context.Context, args interface{}) (interface{}, error)
}

type Result struct {
	value string
	err   error
}

type WorkerPool struct {
	workersCount int
	jobs         chan Job
	results      chan Result
}

func New(workersCount int) WorkerPool {
	return WorkerPool{
		workersCount: workersCount,
		jobs:         make(chan Job, workersCount),
		results:      make(chan Result, workersCount),
	}
}

func (wp WorkerPool) GenerateFrom(jobsBulk []Job) {
	for i, _ := range jobsBulk {
		wp.jobs <- jobsBulk[i]
	}
	close(wp.jobs)
}

func worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan Job, results chan<- Result) {
	defer wg.Done()

	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}
		case <-ctx.Done():
			fmt.Println("cancelled worker. Error detail: %v \n", ctx.Err())
		}
	}
}
