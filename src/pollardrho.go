package main

import "math/big"
import "math/rand"
import "runtime"
import "time"
//~ import "fmt"

type polynomial func(*big.Int) *big.Int

var (
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func get_f(toFactor *big.Int) polynomial {
	return func(x *big.Int) *big.Int {
		r := new(big.Int).Mul(x,x)
		rand_const := rng.Int63()  
		for rand_const < 1 {
			rand_const = rng.Int63()  
		}
		r.Add(r, big.NewInt(rand_const))
		r.Mod(r, toFactor)
		return r
	}
}

func pollardRho(task *Task, toFactor *big.Int, f polynomial) (*big.Int, bool, bool) {
	var x,y,d *big.Int
	x = big.NewInt(2)
	y = big.NewInt(2)
	d = big.NewInt(1)

	for(d.Cmp(big.NewInt(1)) == 0) {	
		if(task.ShouldStop()) {
			return d, false, true
		}
		
		
		x = f(x) 
		y = f(f(y))
		r := new(big.Int)
		r.Sub(x,y)
		r.Abs(r)
		d = r.GCD(nil, nil, r, toFactor)
		    
		// Allow other go threads to run
        runtime.Gosched()
	}
	if(d.Cmp(toFactor) == 0) {
		return d, true, false
	}
	
	return d, false, false
}

func pollardFactoring(task *Task) ([]*big.Int, bool) {	
	//~ fmt.Println("Starting pollard.")
	buffer := make([]*big.Int, 0, 100)
	quo := new(big.Int)
	quo.Set(task.toFactor)
	
	f := get_f(task.toFactor)
	for !quo.ProbablyPrime(prime_precision) {//quo.Cmp(big.NewInt(1)) > 0) {
		if(task.ShouldStop()) {
			return buffer, true
		}
		
		tmp, newQuo, timed_out := trialdivision(task,quo)		
		buffer = append(buffer, tmp...)
		if(timed_out) {
			return buffer, true
		}
		if(newQuo == nil) {
			return buffer, false
		}		
		quo = newQuo		
		

		factor, error, timed_out := pollardRho(task, quo, f)		
		if(timed_out) {
			return buffer, true
		}
		
		if(error || !factor.ProbablyPrime(prime_precision)) {
			// Try again
			f = get_f(task.toFactor)
			continue
		}
		buffer = append(buffer, factor)
        quo.Quo(quo, factor)                                

	}
	buffer = append(buffer, quo)
	return buffer, false
}
