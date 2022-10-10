package main

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
