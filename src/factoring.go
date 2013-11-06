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

type partResult struct {
	index int
	factor *big.Int
}

type factoring func(time.Time, time.Duration, *big.Int) []*big.Int

var (
	allowedRunTime int = 1400 // milliseconds
	numWorkers int
	numTasks int
	inputsize int
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	deadline = 10
	prime_precision = 20
	capacity = 100
	
	resultChannel = make(chan [][]*partResult)
	stopTime = time.Now()
	
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
	activeGoRoutines := 0
	numTasks = len(tasks)
	resultsReceived = 0
	
	// Receive and save results and create new tasks if possible until done
	done := false
	for !done {
		select {
			case result, open := <- resultSubmission:
				if open {
					for _, res := range result {
						if res.factor == nil {
							results[res.index] = nil	
							//~ fmt.Println("Coordinator:", "Failed task", res.index)			
							break
						}
						//~ fmt.Println("Coordinator:", "Received result", res.index)			
						results[res.index] = append(results[res.index], res)
					}
					activeGoRoutines--
					resultsReceived++
				}
			default:
				if activeGoRoutines < numWorkers && nextTask < len(tasks) {				
					t := tasks[nextTask]
					nextTask += 1
					go work(*t)
					activeGoRoutines++					
				} else if nextTask == len(tasks) && resultsReceived == len(tasks) {
					done = true
				} else {
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

func work(task task) {
	// Do task if it is not nil
	rawResult := factorise(&task.toFactor)
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
			break
        } else {
			newTask := new(task)
			newTask.index = int(i)
			newTask.toFactor = *factorValue
			tasks = append(tasks, newTask)
        }
    }
    
    timeout := time.Duration(allowedRunTime) * time.Millisecond
	stopTime = time.Now().Add(timeout)
	//~ fmt.Println("Current time:", time.Now())
	//~ fmt.Println("Timeout set for:", stopTime)
	
	
    inputsize = len(tasks)
    sort.Sort(tasks)
    
    //~ for _, toFactor := range tasks {
		//~ fmt.Println(toFactor)
	//~ }
	

    //~ numProcs := runtime.NumCPU()
    numWorkers = runtime.NumCPU()
	runtime.GOMAXPROCS(numWorkers)
	
	quit := make(chan bool)
	go coordinate(tasks, &quit)
	<- quit
	//~ fmt.Println("Current time:", time.Now())
}
