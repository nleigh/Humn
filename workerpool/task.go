package workerpool

import (
	"fmt"
	"os"
)

type Task struct {
	Err  error
	Data interface{}
	f    func(interface{}) error
}

func NewTask(f func(interface{}) error, data interface{}) *Task {
	return &Task{f: f, Data: data}
}

func process(task *Task) {
	task.Err = task.f(task.Data)
	if task.Err != nil {
		// Any logging performed by the program should be directed to standard error (stderr)
		_, _ = fmt.Fprintln(os.Stderr, task.Err.Error())
	}

}
