package main

import "time"
import "fmt"

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
	for _, task := range tasks {
		duration := stopTime.Sub(time.Now()) / time.Duration(len(tasks) - finishedTasks)
		fmt.Println(duration)
		go task.Run()
		select {
			case <-time.After(duration):
				fmt.Println("Timeout occured.")
				task.Stop()
			case <-task.ch:		
				fmt.Println("Finished normally.")							
		}		
		finishedTasks++
	}
}
