package main

import "fmt"
import "math/big"
import "time"
import "runtime"
import "bufio"
import "os"
import "strings"
import "sort"

type partResult struct {
	index  int
	factor *big.Int
}

type (
    factoring func(*big.Int, chan bool) ([]*big.Int, bool)
	naivefactoring func(*big.Int) ([]*big.Int, *big.Int, bool)
) 


var (
	stopTime time.Time
	resultSubmission chan []*partResult
	
	numWorkers = 1 // Kommer antagligen alltid vara ett för kattis..
	allowedRunTime  int = 14000 // milliseconds
	prime_precision = 20
	resultsReceived = 0
	finishedTasks = 0
	numTasks int
	f factoring = pollardFactoring
)



func appendSlice(thisSlice, toAppend []*big.Int) []*big.Int {
	for _, val := range toAppend {
		thisSlice = append(thisSlice, val)
	}
	return thisSlice
}

// Coordinator main function
// Coordinate task solving and when all is done print results
func coordinate(factoringMethod factoring, tasks Tasks) { 

	// Reinitialize submission channel
	resultSubmission = make(chan []*partResult, len(tasks))

	// Init collection that will hold results
	results := make([][]*partResult, len(tasks))
	for i := 0; i < len(tasks); i++ {
		results[i] = make([]*partResult, 0, len(tasks))
	}

	// Some counters
	nextTask := 0
	activeGoRoutines := 0
	numTasks = len(tasks)
	
	// Receive and save results and create new tasks if possible until done
	done := false
	for !done {
		select {
		case result, open := <-resultSubmission:
			if open {
				for _, res := range result {
					if res.factor == nil {
						results[res.index] = nil							
						break
					}					
					results[res.index] = append(results[res.index], res)
				}
				activeGoRoutines--
				resultsReceived++
			}
		default:
			if done {
				break
			} 
			if nextTask < len(tasks) {
				t := tasks[nextTask]
				nextTask += 1
				activeGoRoutines++
				start_task(t)
				activeGoRoutines--
				
			} else if nextTask == len(tasks) {
				done = true
			} 
			/*
			else if time.Now().Equal(stopTime) || time.Now().After(stopTime) {
				////fmt.Println(duration)fmt.Println("Coordinator:", "Timeout @", t1)
				done = true
			}*/
			//runtime.Gosched()
		}
	}

	printResult(results)
}

// Coordinator function
// Print out the result when all tasks are finished
func printResult(resultCollection [][]*partResult) {
	for _, results := range resultCollection {
		if results == nil || len(results) == 0 {
			fmt.Println("fail")
			fmt.Println("")
			continue
		}
		for _, res := range results {
			if res.factor.Cmp(big.NewInt(0)) == 0 || res.factor.Cmp(big.NewInt(1)) == 0 {
				fmt.Println("fail")
				break
			}
			fmt.Println(res.factor)
		}
		fmt.Println("")
	}
}

func work(task *Task) {	
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
	res, timed_out := f(task.toFactor, task.ch)			
	if(timed_out) {
		//~ fmt.Println("Timeout scoped out")
		return
	}					
	// rawResult := appendSlice(rawResult, res)		
	doResultSubmission(task.index, res)
}

// Worker main function
// Do work with task and submit answer through global resultSubmission (channel)
func start_task(task *Task) {

	//rawResult := make([]*big.Int, 0, 15)
	// newFactor := new(big.Int).Set(task.toFactor)

	duration := stopTime.Sub(time.Now()) / time.Duration(numTasks - finishedTasks)
	//~ fmt.Println(duration)
	go task.Run()
	select {
		case <-time.After(duration):
			//~ fmt.Println("Timeout occured.")
			task.Stop()
		case <-task.ch:		
			//~ fmt.Println("Finished normally.")							
	}
	//~ fmt.Println("Finished task")
	
	finishedTasks++	

	return
}

func doResultSubmission(taskId int, rawResult []*big.Int) {
	result := make([]*partResult, 0, 100)
	if rawResult == nil {
		res := partResult{taskId, nil}
		result = append(result, &res)
		//////fmt.Println("Worker:", "Exeeded time limit of", timeout)
	} else {
		for _, rawRes := range rawResult {
			res := partResult{taskId, rawRes}
			result = append(result, &res)
		}
	}
	// Send to coordinator
	resultSubmission <- result
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
	////fmt.Println("Timeout is", timeout)
	stopTime = time.Now().Add(timeout)
	sort.Sort(tasks)

	//~ for _, toFactor := range tasks {
	//~ ////fmt.Println(toFactor)
	//~ }

	runtime.GOMAXPROCS(numWorkers)

	//quit := make(chan bool, 1)
	coordinate(pollardFactoring, tasks)
	//<-quit	
	////fmt.Println("Time elapsed", time.Now().Sub(start))
}
