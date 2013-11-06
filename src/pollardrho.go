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

func pollardRho(toFactor *big.Int, f polynomial, finished chan bool) (*big.Int, bool, bool) {
	var x,y,d *big.Int
	x = big.NewInt(2)
	y = big.NewInt(2)
	d = big.NewInt(1)

	for(d.Cmp(big.NewInt(1)) == 0) {	
		select {
			case <-finished:			
				//~ fmt.Println("Timeout signal 2!")				
				return d, false, true
			default:
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

func pollardFactoring(toFactor *big.Int, finished chan bool) ([]*big.Int, bool) {	
	//~ fmt.Println("Starting pollard.")
	buffer := make([]*big.Int, 0, 100)
	quo := new(big.Int)
	quo.Set(toFactor)
	
	f := get_f(toFactor)
	for !quo.ProbablyPrime(prime_precision) {//quo.Cmp(big.NewInt(1)) > 0) {
		select {
			case <-finished:
				//~ fmt.Println("Timeout signal 3!")
				return buffer, true
			default:
		}
		
		tmp, newQuo, timed_out := trialdivision(quo, finished)		
		buffer = append(buffer, tmp...)
		if(timed_out) {
			return buffer, true
		}
		if(newQuo == nil) {
			return buffer, false
		}		
		quo = newQuo		
		

		factor, error, timed_out := pollardRho(quo, f, finished)		
		if(timed_out) {
			return buffer, true
		}
		
		if(error || !factor.ProbablyPrime(prime_precision)) {
			// Try again
			f = get_f(toFactor)
			continue
		}
		buffer = append(buffer, factor)
        quo.Quo(quo, factor)                                

	}
	buffer = append(buffer, quo)
	return buffer, false
}
