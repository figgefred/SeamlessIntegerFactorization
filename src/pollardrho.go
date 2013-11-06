package main

import "math/big"
import "math/rand"
import "runtime"
import "time"

type polynomial func(*big.Int) *big.Int

var (
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func get_f(toFactor *big.Int) polynomial {
	return func(x *big.Int) *big.Int {
		r := new(big.Int).Mul(x,x)
		r.Add(r, big.NewInt(rng.Int63()))
		r.Mod(r, toFactor)
		return r
	}
}

func pollardRho(toFactor *big.Int, f polynomial) (*big.Int, bool) {
	var x,y,d *big.Int
	x = big.NewInt(2)
	y = big.NewInt(2)
	d = big.NewInt(1)


	for(d.Cmp(big.NewInt(1)) == 0) {
	
		x = f(x) 
		y = f(f(y))
		//~ ////fmt.Println(x)
		//~ ////fmt.Println(y)
		//~ ////fmt.Println()
		r := new(big.Int)
		r.Sub(x,y)
		r.Abs(r)
		d = r.GCD(nil, nil, r, toFactor)
    
		// Allow other go threads to run
        runtime.Gosched()

	}
	if(d.Cmp(toFactor) == 0) {
		return d, true
	}
	
	return d, false
}

func pollardFactoring(toFactor *big.Int) []*big.Int {	
	buffer := make([]*big.Int, 0, 100)
	
	quo := new(big.Int)
	quo.Set(toFactor)
	
	for !quo.ProbablyPrime(prime_precision) {//quo.Cmp(big.NewInt(1)) > 0) {

		f := get_f(toFactor)
		factor, error := pollardRho(quo, f)
		
		if(error || !factor.ProbablyPrime(prime_precision)) {
			// Try again
			continue
		}
		buffer = append(buffer, factor)
        quo.Quo(quo, factor)                                

	}
	buffer = append(buffer, quo)
	return buffer
}