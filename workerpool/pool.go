package workerpool

import (
	"sync"
	"time"
)

type Pool struct {
	Tasks   []*Task
	Workers []*Worker

	concurrency   int
	collector     chan *Task
	runBackground chan bool
	wg            sync.WaitGroup
}

// NewPool initializes a new pool with the given tasks and
// at the given concurrency.
func NewPool(tasks []*Task, concurrency int) *Pool {
	return &Pool{
		Tasks:       tasks,
		concurrency: concurrency,
		collector:   make(chan *Task, 1000),
	}
}

// AddTask adds a task to the pool
func (p *Pool) AddTask(task *Task) {
	p.collector <- task
}

// RunBackground runs the pool in background
func (p *Pool) RunBackground() {

	for i := 1; i <= p.concurrency; i++ {
		worker := NewWorker(p.collector, i)
		p.Workers = append(p.Workers, worker)
		go worker.StartBackground()
	}

	for i := range p.Tasks {
		p.collector <- p.Tasks[i]
	}

	p.runBackground = make(chan bool)
	<-p.runBackground
}

// Stop stops background workers
func (p *Pool) Stop() {
	for i := range p.Workers {

		// wait for workers to finish before stopping
		for p.Workers[i].Working == true {
			time.Sleep(100 * time.Millisecond)
		}
		p.Workers[i].Stop()
	}
	p.runBackground <- true
}
