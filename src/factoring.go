package main

import "fmt"
import "math/big"
import "math/rand"
import "time"
import "runtime"
import "bufio"
import "os"
import "strings"
import "sort"
//import "strconv"

type partResult struct {
	index int
	factor *big.Int
}

type factoring func(time.Time, time.Duration, *big.Int) []*big.Int

var (
	totalRunTime int = 15000 // milliseconds

	numWorkers int
	inputsize int
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	deadline = 10
	prime_precision = 20
	capacity = 100
	
	resultChannel = make(chan [][]*partResult)

	resultSubmission chan []*partResult
	taskChannel chan *task
)

// A function that should be runned in a go routine.
// This function initializes new worker Go-routines 
// for every task.
// It will only do NumProcs workers at a time.
func coordinate(tasks Tasks, finishedChan *chan bool) {

	// Reinitialize submission channel
	resultSubmission = make(chan []*partResult, numWorkers)

	// Init collection that will hold results
	results := make([][]*partResult, inputsize)
	for i:= 0; i < inputsize; i++ {
		results[i] = make([]*partResult, 0, capacity)
	}
	
	// Some counters
	nextTask := 0
	resultsReceived := 0
	runTime := totalRunTime
	activeGoRoutines := 0

	// Receive and save results and create new tasks if possible until done
	done := false
	for !done {
		select {
			case result, open := <- resultSubmission:
				if open {
					for _, res := range result {
						if res.factor == nil {
							results[res.index] = nil	
							//fmt.Println("Coordinator:", "Failed task", res.index)			
							//fmt.Println("Coordinator:", "Trashed task", res.index)			
							break
						}
						////fmt.Println("Coordinator:", "Received result", res.index)			
						results[res.index] = append(results[res.index], res)
					}
					activeGoRoutines--
					resultsReceived++
				}
			default:
				if activeGoRoutines < numWorkers && nextTask < len(tasks) {
					for activeGoRoutines+1 <= numWorkers {
						if nextTask == len(tasks) {
							break
						}
						t := tasks[nextTask]
						nextTask += 1
						runTime = (totalRunTime/len(tasks))
						timeout := time.Duration(runTime)*time.Millisecond
						//fmt.Println("Coordinator:", "Setting timelimit @", timeout, "for Worker", nextTask)
						start := time.Now()
						go work(start, timeout, *t, pollardFactoring)
						activeGoRoutines++
					}
				} else if nextTask == len(tasks) && resultsReceived == len(tasks) {
					////fmt.Println("Coordinator:", "Now we are done")
					done = true
				} else {
					////fmt.Println("Coordinator:", "Nothing to do")
					runtime.Gosched()
				}
		}
	}
	close(resultSubmission)
	//elapsedTime := time.Since(initTime)
	// Dump out to sys out
	printResult(results)
	////fmt.Println("Coordinator:", "Finished after", elapsedTime)
	*finishedChan <- true
}

func work(start time.Time, timeout time.Duration, task task, f factoring) {
	
	// Do task if it is not nil
	rawResult := f(start, timeout, &task.toFactor)
	result := make([]*partResult, 0, 100)
	if rawResult == nil {
		res := partResult{task.index, nil}
		result = append(result, &res)
	//fmt.Println("Worker:", "Exeeded time limit of", timeout)
	}
	for _, rawRes := range rawResult {
		res := partResult{task.index, rawRes}
		result = append(result, &res)
	}
	// Send to coordinator
	resultSubmission <- result
}

func printResult(resultCollection [][]*partResult) {
	for _, results := range resultCollection {
		if results == nil {
			fmt.Println("fail")
			fmt.Println("")
			continue
		}
		for _, res := range results {
			if res.factor.Cmp(big.NewInt(0)) == 0 {
				fmt.Println("fail")
				break
			}
			fmt.Println(res.factor)
		}
		fmt.Println("")
	}
}

func main() {
		
	reader := bufio.NewReader(os.Stdin)

	factorCount := 100
	tasks := make(Tasks, 0, factorCount)
	// Read in line by line
	for i:=0; ; i++ {
        line, _ := reader.ReadString('\n')
    	if(strings.TrimSpace(line) == "") {
    		break
    	}

        factorValue, ok := (new(big.Int)).SetString(strings.TrimSpace(line), 10) 
        if !ok {
        	////fmt.Println("Parse error of", line)
			// Exit
        	return
        } else {
        	newTask := new(task)
        	newTask.index = int(i)
        	newTask.toFactor = *factorValue
        	tasks = append(tasks, newTask)
        }
    }
    inputsize = len(tasks)
    sort.Sort(tasks)

    //for _, t  := range tasks{
    	//fmt.Println(t.toFactor)
    //}

    //numProcs := runtime.NumCPU()
    numProcs := 1
	runtime.GOMAXPROCS(numProcs)

	numWorkers = numProcs
	
	quit := make(chan bool)
	go coordinate(tasks, &quit)
	<- quit
	/*if( numProcs > 1) {
		for _, _ = range factorValues {
			////fmt.Println("fail")
			////fmt.Println("")
		}
		return
	}*/
	/*
	for _, toFactor := range factorValues {

		//timeout := time.After(time.Duration(deadline) * time.Second)
		resultChannel = make(chan string)
		timeout := make(chan bool, 1)
		
		go func() {
			resultChannel <- factorise(toFactor)
		}();
		
		go func() {
			time.Sleep(time.Duration(deadline) * time.Millisecond)	
			timeout <- true
		}();
		
		select {		
			case factors := <- resultChannel:
				////fmt.Println(factors)	
			case <- timeout:
				////fmt.Println("fail")	
				////fmt.Println()		
		}
	}*/
}
