package main

import "math/big"
import "runtime"

type polynomial func(*big.Int) *big.Int

func pollardRho(task *Task, toFactor *big.Int) (*big.Int, bool) {	
	/*
	x := TWO
	y := TWO
	*/
	x := new(big.Int).Rand(rng, toFactor)
	y := new(big.Int).Rand(rng, toFactor)
	d := ONE	
	r := new(big.Int)
	rand_const := new(big.Int).Rand(rng, toFactor)
	
	if(r.Mod(toFactor, TWO).Cmp(ZERO) == 0) {
		return TWO, false
	}
	
	i := 0
	// i < 10 to prevent a bad random seed from finding a factor.
	for(d.Cmp(ONE) == 0 && i < 10) { 
		i++	
		if(task.ShouldStop()) {
			return d, false
		}		
		
		x = x.Mul(x,x).Add(x, rand_const).Mod(x, toFactor)
		y = y.Mul(y,y).Add(y, rand_const).Mod(y, toFactor)
		y = y.Mul(y,y).Add(y, rand_const).Mod(y, toFactor)
		//~ y = f(f(y))
		r.Sub(x,y).Abs(r)		
		d = r.GCD(nil, nil, r, toFactor)
	}
	if(d.Cmp(toFactor) == 0) {
		return d, true
	}
	
	return d, false
}

func pollardFactoring(task *Task) ([]*big.Int) {	
	return _pollardFactoring(task, task.toFactor)
}

func _pollardFactoring(task *Task, toFactor *big.Int) ([]*big.Int) {
	buffer := make([]*big.Int, 0, 100)
	quo := new(big.Int)
	quo.Set(task.toFactor)
	
	//~ f := get_f(task.toFactor)
	for !quo.ProbablyPrime(prime_precision) {//quo.Cmp(big.NewInt(1)) > 0) {
		if(task.ShouldStop()) {
			return buffer
		}
		
		/*
		tmp, newQuo, timed_out := trialdivision(task,quo)		
		buffer = append(buffer, tmp...)
		if(timed_out) {
			return buffer, true
		}
		if(newQuo == nil) {
			return buffer, false
		}		
		quo = newQuo		
		*/

		factor, error := pollardRho(task, quo)	
		
		if(error || !factor.ProbablyPrime(prime_precision)) {
			// Allow other go threads to run
			runtime.Gosched() 
			// Try again
			//~ f = get_f(task.toFactor)
			continue
		}
		
		
		if(!factor.ProbablyPrime(prime_precision)) {
			sub := _pollardFactoring(task, factor)
			buffer = append(buffer, sub...)		
		}
	
		buffer = append(buffer, factor)
        quo.Quo(quo, factor)    
        
                          
	}
	return append(buffer, quo)
}
