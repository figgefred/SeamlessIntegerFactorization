package main

import "time"
import "sync"

type Tasks []*Task

func (tasks Tasks) Len() int {
	return len(tasks)
}

func (tasks Tasks) Less(i, j int) bool {
	return (tasks[i].toFactor).Cmp(tasks[j].toFactor) == -1 
}

func (tasks Tasks) Swap(i, j int) {
	tasks[i], tasks[j] = tasks[j], tasks[i]
}

func (tasks Tasks) PrintResults() {
	for _, task := range tasks {
		task.PrintResult()
	}
}

func (tasks Tasks) RunTasksWithTimeout(stopTime time.Time) {		
	finishedTasks := 0
	waitCond := new(sync.Cond)
	waitCond.L = new(sync.Mutex)

	done := make(chan bool, 1)
	current_running := 0
	for _, task := range tasks {				
		
		waitCond.L.Lock()
		for(current_running > numWorkers) {
			waitCond.Wait()
		}
		current_running++
		waitCond.L.Unlock()
		
		go func(task* Task) {					
			todo := len(tasks) - finishedTasks - numWorkers
			if(todo < 1) {
				todo = 1
			}
			duration := stopTime.Sub(time.Now()) / time.Duration(todo)
			dprint(duration)
			go task.Run()
			select {
				case <-time.After(duration):
					dprint("Timeout occured.")
					task.Stop()
				case <-task.ch:		
					dprint("Finished normally.")							
			}		
				
			waitCond.L.Lock()	
			current_running--
			finishedTasks++
			if(finishedTasks == len(tasks)) {
				done <- true
			}	
			waitCond.Signal()
			waitCond.L.Unlock()	
		}(task)
	}
	<-done
}
