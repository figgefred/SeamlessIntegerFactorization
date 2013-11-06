package main

import "math/big"

type task struct {
	index int
	toFactor big.Int
}

type Tasks []*task

func (tasks Tasks) Len() int {
	return len(tasks)
}

func (tasks Tasks) Less(i, j int) bool {
	return (&tasks[i].toFactor).Cmp(&tasks[j].toFactor) == -1 
}

func (tasks Tasks) Swap(i, j int) {
	tasks[i], tasks[j] = tasks[j], tasks[i]
}