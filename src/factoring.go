package main

import "fmt"
import "math/big"

func gcd(a big.Int, b big.Int) float64 {
	if(b == 0) {
		return a
	}
	if(a > b) {
		return gcd(a-b, b)
	} else if(a < b) {
		return gcd(a,b-a)
	} else {
		return a
	}
	
	return 1
}

type polynomial func(big.Int) big.Int

func pollardRho(toFactor big.Int, f polynomial) (big.Int, bool) {
	var x,y,d big.Int
	x = 2
	y = 2
	d = 1
	for(d == 1) {
		x = f(x) 
		y = f(f(y))
		d = gcd(big.Abs(x-y),toFactor)
	}
	if(d == toFactor) {
		return d, true
	}
	
	return d, false
}

func get_f(toFactor big.Int) polynomial {
	return func(x big.Int) big.Int {
		return big.Mod(x*x - 1, toFactor)
	}
}

func factorise(toFactor float64) {
	// factors := make([]int, 20)
		
}

func main() {
	numberToFactor = new(big.Int)
	fmt.Scan(&toFactor)

	fmt.Println(pollardRho(toFactor,  get_f(toFactor)))
}
