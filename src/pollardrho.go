package main

import "math/big"

type polynomial func(*big.Int) *big.Int

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
	}
	if(d.Cmp(toFactor) == 0) {
		return d, true
	}
	
	return d, false
}

func pollardFactoring(toFactor *big.Int) []*big.Int {	
	buffer := make([]*big.Int, 0, 100)
	if(toFactor.ProbablyPrime(prime_precision)) {
		return append(buffer, toFactor)
	}
	
	quo := new(big.Int)
	quo.Set(toFactor)
	
	for(quo.Cmp(big.NewInt(1)) > 0) {

		f := get_f(toFactor)
		factor, error := pollardRho(quo, f)
		
		if(error || factor.Int64() == int64(0)) {
			// Try again
			continue
		}
		
        quo.Quo(quo, factor)                                
        
        if(!factor.ProbablyPrime(prime_precision)) {	
        	res := pollardFactoring(factor)
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
	}
	// Lets redo this - send back old task for hope of better function
	return buffer
}