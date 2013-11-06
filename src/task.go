package main

import "math/big"
import "sync"
//~ import "fmt"

type doWork func(*Task) 

type Task struct {
	index int
	toFactor *big.Int
	
	ch        chan bool
	waitGroup *sync.WaitGroup
	finished bool
	w doWork
}

type Tasks []*Task

// Make a new Task.
func NewTask(index int, toFactor *big.Int, w doWork) *Task {
	t := &Task {
		ch: make(chan bool),
		waitGroup: &sync.WaitGroup{},		
		index: index,
		toFactor: toFactor,
		finished: false,
		w: w,
	}
	t.waitGroup.Add(1)
	return t
}

func (task* Task) Stop() {	
	//~ close(task.ch)
	task.ch <- true
	task.finished = true
	task.waitGroup.Wait()
	//~ fmt.Println("no longer waiting.")
}

func (task* Task) ShouldStop() bool {
	if task.finished {
		return task.finished
	}
	
	select {
		case <-task.ch:
			task.finished = true
			return true
		default:			
	}
	return false
}

func (task* Task) Run() {
	defer task.waitGroup.Done()
	defer close(task.ch)
	task.w(task)
	//~ task.ch <- true
}

func (tasks Tasks) Len() int {
	return len(tasks)
}

func (tasks Tasks) Less(i, j int) bool {
	return (tasks[i].toFactor).Cmp(tasks[j].toFactor) == -1 
}

func (tasks Tasks) Swap(i, j int) {
	tasks[i], tasks[j] = tasks[j], tasks[i]
}
