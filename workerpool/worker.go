package workerpool

// Worker handles all the work
type Worker struct {
	ID       int
	taskChan chan *Task
	quit     chan bool
	Working  bool
}

// NewWorker returns new instance of worker
func NewWorker(channel chan *Task, ID int) *Worker {
	return &Worker{
		ID:       ID,
		taskChan: channel,
		quit:     make(chan bool),
	}
}

// StartBackground starts the worker in background waiting
// Added boolean flag for if worker is working
func (wr *Worker) StartBackground() {

	for {
		select {
		case task := <-wr.taskChan:
			{
				wr.Working = true
				process(task)
				wr.Working = false
			}
		case <-wr.quit:
			return
		}
	}
}

// Stop quits the worker
func (wr *Worker) Stop() {
	go func() {
		wr.quit <- true
	}()
}
