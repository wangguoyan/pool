package pool

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

type Task struct {
	f    func()
	name string
}

func newTask(f func()) *Task {
	t := Task{
		f: f,
	}
	return &t
}

func (t *Task) exec() {
	t.f()
}

type Pool struct {
	entryChannel chan *Task

	jobsChannel chan *Task

	workerNum int

	stopChan chan bool
}

func (p *Pool) AddTask(f func()) *Task {
	t := newTask(f)
	p.entryChannel <- t
	return t
}
func NewPool(cap int) *Pool {
	return &Pool{
		entryChannel: make(chan *Task),
		workerNum:    cap,
		jobsChannel:  make(chan *Task),
		stopChan:     make(chan bool),
	}
}

func (p *Pool) worker() {
	for {
		select {
		case task := <-p.jobsChannel:
			var funcName string
			now := time.Now().UnixNano() / 1e6
			task.exec()
			if task.name != "" {
				funcName = task.name
			} else {
				funcName = runtime.FuncForPC(reflect.ValueOf(task.f).Pointer()).Name()
			}
			methodTime := (time.Now().UnixNano() / 1e6) - now
			msg := fmt.Sprintf("exec task %30s - time: %vms", funcName, methodTime)
			fmt.Println(msg)
		case <-p.stopChan:
			close(p.jobsChannel)
			return
		}
	}
}

func (p *Pool) Run() *Pool {
	for i := 0; i < p.workerNum; i++ {
		go p.worker()
	}
	go func() {
		for {
			select {
			case task := <-p.entryChannel:
				p.jobsChannel <- task
			case <-p.stopChan:
				close(p.entryChannel)
				return
			}
		}

	}()
	return p
}

func (p *Pool) Stop() {
	p.stopChan <- true
}
