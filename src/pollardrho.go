package main

import "math/big"

func pollardRho(task *Task, toFactor *big.Int) (*big.Int, bool) {	
	/*
	x := TWO
	y := TWO
	*/
	
	x := new(big.Int).Rand(rng, toFactor)
	y := new(big.Int).Rand(rng, toFactor)
	
	d := ONE	
	r := new(big.Int)
	rand_const := r.Rand(rng, toFactor)
	
	if(r.Mod(toFactor, TWO).Cmp(ZERO) == 0) {
		return TWO, false
	}
	
	i := 0
	// i < 10 to prevent a bad random seed from finding a factor.
	for(d.Cmp(ONE) == 0 && i < 10) { 
		i++	
		/*
		if(task.ShouldStop()) {
			return d, false
		}*/		
		
		x = x.Mul(x,x).Add(x, rand_const).Mod(x, toFactor)
		y = y.Mul(y,y).Add(y, rand_const).Mod(y, toFactor)
		y = y.Mul(y,y).Add(y, rand_const).Mod(y, toFactor)		
		r.Sub(x,y).Abs(r)		
		d = r.GCD(nil, nil, r, toFactor)
	}
	if(i == 10 || d.Cmp(toFactor) == 0) {
		return d, true
	}
	
	return d, false
}

func pollardFactoring(task *Task) ([]*big.Int) {	
	return _pollardFactoring(task, task.toFactor)
}

func _pollardFactoring(task *Task, toFactor *big.Int) ([]*big.Int) {
	buffer := make([]*big.Int, 0)
	quo := new(big.Int)
	quo.Set(toFactor)
	
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
		
		if(error || factor.Cmp(ONE) == 0) {
			continue
		}

		if(!factor.ProbablyPrime(prime_precision)) {
			sub := _pollardFactoring(task, factor)
			buffer = append(buffer, sub...)		
		} else {	
			buffer = append(buffer, factor)
		}     
		   
        quo.Quo(quo, factor)         
	}
	return append(buffer, quo)
}
