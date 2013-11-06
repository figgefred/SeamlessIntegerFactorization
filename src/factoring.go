package main

//~ import "fmt"
import "math/big"
import "time"
import "runtime"
import "bufio"
import "os"
import "strings"
import "sort"

type (
    factoring func(*big.Int, chan bool) ([]*big.Int, bool)
	naivefactoring func(*big.Int) ([]*big.Int, *big.Int, bool)
) 

var (
	numWorkers = 1 // Kommer antagligen alltid vara ett för kattis..
	allowedRunTime  int = 14000 // milliseconds
	prime_precision = 20

	f factoring = pollardFactoring
)

func work(task *Task) ([]*big.Int, bool) {	
	//rawResult, newFactor, timed_out = trialdivision(newFactor, task.ch)		
	/*
	if(timed_out) {
		return
	}		
	
	// We are done
	if newFactor == nil {
		doResultSubmission(task.index, rawResult)
		task.ch <- true
		return
	}
	*/
	
	// Do expensive factorization
	return f(task.toFactor, task.ch)				
	
						
	// rawResult := appendSlice(rawResult, res)		
	
}

func main() {

	reader := bufio.NewReader(os.Stdin)
	//start := time.Now()
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
			newTask := NewTask(int(i), factorValue, work)			
			tasks = append(tasks, newTask)
		}
	}

	timeout := time.Duration(allowedRunTime) * time.Millisecond	
	stopTime := time.Now().Add(timeout)
	
	sort.Sort(tasks)

	runtime.GOMAXPROCS(numWorkers)

	tasks.RunTasksWithTimeout(stopTime)
	tasks.PrintResults()
}
