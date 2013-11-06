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

type task struct {
	index int
	toFactor big.Int
}

type partResult struct {
	index int
	factor *big.Int
}

type polynomial func(*big.Int) *big.Int
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
						go work(start, timeout, *t)
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

func work(start time.Time, timeout time.Duration, task task) {
	
	// Do task if it is not nil
	rawResult := factorise(start, timeout, &task.toFactor)
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

func pollardRho(start time.Time, timeout time.Duration, toFactor *big.Int, f polynomial) (*big.Int, bool) {
	var x,y,d *big.Int
	x = big.NewInt(2)
	y = big.NewInt(2)
	d = big.NewInt(1)
	currTime := time.Now()
	for(currTime.Sub(start) < timeout && d.Cmp(big.NewInt(1)) == 0) {
		x = f(x) 
		y = f(f(y))
		//~ ////fmt.Println(x)
		//~ ////fmt.Println(y)
		//~ ////fmt.Println()
		r := new(big.Int)
		r.Sub(x,y)
		r.Abs(r)
		d = r.GCD(nil, nil, r, toFactor)
		currTime = time.Now()
	}
	if currTime.Sub(start) > timeout {
		return nil, true
	}
	if(d.Cmp(toFactor) == 0) {
		return d, true
	}
	
	return d, false
}

func get_f(toFactor *big.Int) polynomial {
	return func(x *big.Int) *big.Int {
		r := new(big.Int).Mul(x,x)
		r.Add(r, big.NewInt(rng.Int63()))
		r.Mod(r, toFactor)
		return r
	}
}

func factorise(start time.Time, timeout time.Duration, toFactor *big.Int) []*big.Int {	
	buffer := make([]*big.Int, 0, 100)
	if(toFactor.ProbablyPrime(prime_precision)) {
		return append(buffer, toFactor)
	}
	
	quo := new(big.Int)
	quo.Set(toFactor)
	currTime := time.Now()
	for(currTime.Sub(start) < timeout && quo.Cmp(big.NewInt(1)) > 0) {

		f := get_f(toFactor)
		factor, error := pollardRho(start, timeout, quo, f)
		
		if(error || factor.Int64() == int64(0)) {
			// Try again
			currTime = time.Now()
			continue
		}
		
        quo.Quo(quo, factor)                                
        
        if(!factor.ProbablyPrime(prime_precision)) {	
        	res := factorise(start, timeout, factor)
        	if res == nil {
        		return nil
        	}
        	for _, r := range res {
        		buffer = append(buffer, r)
        	}
        } else {
        	buffer = append(buffer, factor)
        }

        if(quo.ProbablyPrime(prime_precision)) {
            buffer = append(buffer, quo)
            break
        }
        currTime = time.Now()
	}
	if currTime.Sub(start) > timeout {
		return nil
	}
	// Lets redo this - send back old task for hope of better function
	return buffer
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
