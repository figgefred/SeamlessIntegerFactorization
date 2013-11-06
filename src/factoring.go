package main

import "fmt"
import "math/big"
import "math/rand"
import "time"
import "runtime"
import "bufio"
import "os"
import "strings"
//import "sort"

type partResult struct {
	index  int
	factor *big.Int
}

type factoring func(*big.Int) []*big.Int
type naivefactoring func(*big.Int) ([]*big.Int, *big.Int)

var (
	allowedRunTime  int = 14000 // milliseconds
	numWorkers      int
	numTasks        int
	inputsize       int
	rng             = rand.New(rand.NewSource(time.Now().UnixNano()))
	deadline        = 10
	prime_precision = 20
	capacity        = 100

	resultChannel = make(chan [][]*partResult)
	stopTime      time.Time

	resultSubmission chan []*partResult
	taskChannel      chan *Task
)

// Coordinator main function
// Coordinate task solving and when all is done print results
func coordinate(factoringMethod factoring, tasks Tasks, finishedChan *chan bool) {

	// Reinitialize submission channel
	resultSubmission = make(chan []*partResult, len(tasks))

	// Init collection that will hold results
	results := make([][]*partResult, inputsize)
	for i := 0; i < inputsize; i++ {
		results[i] = make([]*partResult, 0, capacity)
	}

	// Some counters
	nextTask := 0
	resultsReceived := 0
	activeGoRoutines := 0
	numTasks = len(tasks)
	resultsReceived = 0

	// Receive and save results and create new tasks if possible until done
	done := false
	for !done {
		select {
		case result, open := <-resultSubmission:
			if open {
				for _, res := range result {
					if res.factor == nil {
						results[res.index] = nil
						//~ //fmt.Println("Coordinator:", "Failed task", res.index)			
						break
					}
					//~ //fmt.Println("Coordinator:", "Received result", res.index)			
					results[res.index] = append(results[res.index], res)
				}
				activeGoRoutines--
				resultsReceived++
			}
		default:
			if done {
				break
			} else if activeGoRoutines < numWorkers && nextTask < len(tasks) {
				t := tasks[nextTask]
				nextTask += 1
				go work(*t, factoringMethod)
				activeGoRoutines++
			} else if nextTask == len(tasks) && resultsReceived == len(tasks) {
				//fmt.Println("Coordinator:", "Finished work @", t1)
				done = true
			} else if time.Now().Equal(stopTime) || time.Now().After(stopTime) {
				//fmt.Println("Coordinator:", "Timeout @", t1)
				done = true
			}
			runtime.Gosched()
		}
	}
	//fmt.Println("Coordinator:", "Done")
	//elapsedTime := time.Since(initTime)
	// Dump out to sys out
	printResult(results)
	//////fmt.Println("Coordinator:", "Finished after", elapsedTime)
	*finishedChan <- true
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

// Worker main function
// Do work with task and submit answer through global resultSubmission (channel)
func work(task Task, f factoring) {

	// Lets try to shorten the value
	rawResult, newFactor := trialdivision(task.toFactor)

	// We are done
	if newFactor == nil {
		doResultSubmission(task.index, rawResult)
		return
	}

	// Do expensive factorization
	res := f(newFactor)
	rawResult = append(rawResult)
	for _, r := range res {
		rawResult = append(rawResult, r)
	}
	doResultSubmission(task.index, rawResult)
	return
}

func doResultSubmission(taskId int, rawResult []*big.Int) {
	result := make([]*partResult, 0, 100)
	if rawResult == nil {
		res := partResult{taskId, nil}
		result = append(result, &res)
		////fmt.Println("Worker:", "Exeeded time limit of", timeout)
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
			newTask := new(Task)
			newTask.index = int(i)
			newTask.toFactor = *factorValue
			tasks = append(tasks, newTask)
		}
	}

	timeout := time.Duration(allowedRunTime) * time.Millisecond
	//fmt.Println("Timeout is", timeout)
	stopTime = time.Now().Add(timeout)

	inputsize = len(tasks)
	//sort.Sort(tasks)

	//~ for _, toFactor := range tasks {
	//~ //fmt.Println(toFactor)
	//~ }

	//numWorkers = runtime.NumCPU()
	numWorkers = 1
	runtime.GOMAXPROCS(numWorkers)

	quit := make(chan bool, 1)
	go coordinate(pollardFactoring, tasks, &quit)
	<-quit	
	//fmt.Println("Time elapsed", time.Now().Sub(start))
}
