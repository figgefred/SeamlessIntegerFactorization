package main

import "fmt"
import "math/big"
import "time"
import "runtime"
import "bufio"
import "os"
import "strings"
import "sort"

var (
	numWorkers = 1 //runtime.NumCPU() // Kommer antagligen alltid vara ett för kattis..
	allowedRunTime  int = 10000 // milliseconds
	work_function = pollardFactoring
	//work_function = naivefactoring
	debug = false
)

func dprint(a ...interface{}) {
	if(!debug) {
		return
	}
	fmt.Println(a...)
}


func main() {	
	dprint("[DEBUG] is on!!!")
	reader := bufio.NewReader(os.Stdin)
	factorCount := 100
	tasks := make(Tasks, 0, factorCount)
	start := time.Now()
	// Read in line by line
	for i := 0; ; i++ {
		line, err := reader.ReadString('\n')
		if err != nil || strings.TrimSpace(line) == "" {
			break
		}

		factorValue, ok := (new(big.Int)).SetString(strings.TrimSpace(line), 10)
		if !ok {
			break
		} else {
			newTask := NewTask(int(i), factorValue, work_function)			
			tasks = append(tasks, newTask)
		}
	}

	timeout := time.Duration(allowedRunTime) * time.Millisecond	
	stopTime := start.Add(timeout)
	
	// Remember to sort results if you turn this on again.
	sort.Sort(tasks)

	runtime.GOMAXPROCS(numWorkers)

	tasks.RunTasksWithTimeout(stopTime)
	executed := time.Since(start)
	
	//~ // Make sure that results are in right order
	results := make(Tasks, len(tasks))
	for _, task := range tasks {
		results[task.index] = task
	}
	results.PrintResults()
	dprint("Executed", executed)
}
