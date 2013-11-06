package main

import "math/big"
import "time"

type polynomial func(*big.Int) *big.Int

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

func pollardFactoring(start time.Time, timeout time.Duration, toFactor *big.Int) []*big.Int {	
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
        	res := pollardFactoring(start, timeout, factor)
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