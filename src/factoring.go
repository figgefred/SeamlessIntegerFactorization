package main

import "fmt"
import "math/big"
import "time"
import "runtime"
import "bufio"
import "os"
import "strings"
//~ import "sort"

var (
	numWorkers = 1 // Kommer antagligen alltid vara ett för kattis..
	allowedRunTime  int = 14500 // milliseconds
	prime_precision = 20
	debug = true
)

func dprint(a ...interface{}) {
	if(!debug) {
		return
	}
	fmt.Println(a...)
}


func main() {	
	dprint("[DEBUG] is on!")
	reader := bufio.NewReader(os.Stdin)
	factorCount := 100
	tasks := make(Tasks, 0, factorCount)
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
			newTask := NewTask(int(i), factorValue, pollardFactoring)			
			tasks = append(tasks, newTask)
		}
	}

	timeout := time.Duration(allowedRunTime) * time.Millisecond	
	stopTime := time.Now().Add(timeout)
	
	// Remember to sort results if you turn this on again.
	//~ sort.Sort(tasks)

	runtime.GOMAXPROCS(numWorkers)

	tasks.RunTasksWithTimeout(stopTime)
	tasks.PrintResults()
}
