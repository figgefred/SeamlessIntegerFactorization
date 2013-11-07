package main

import "math/big"
import "sync"
import "fmt"

type doWork func(*Task) ([]*big.Int)

type Task struct {
	index int
	toFactor *big.Int
	
	ch 	chan bool
	waitGroup *sync.WaitGroup
	finished bool
	timed_out bool
	w doWork
	result []*big.Int 
}

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
	task.timed_out = true
	task.waitGroup.Wait()
}

func (task* Task) PrintResult() {
	if task.timed_out {
		fmt.Println("fail")
		fmt.Println("")
		return
	} 
	
	for _, res := range task.result {
		fmt.Println(res)
	}
	fmt.Println("")
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
	task.result = task.w(task)	
}

func (task* Task) setResults(result []*big.Int) {
	task.result = result
}


