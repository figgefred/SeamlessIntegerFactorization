package main

import "fmt"
import "math/big"
import "math/rand"
import "time"

type polynomial func(*big.Int) *big.Int

var (
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func pollardRho(toFactor *big.Int, f polynomial) (*big.Int, bool) {
	var x,y,d *big.Int
	x = big.NewInt(2)
	y = big.NewInt(2)
	d = big.NewInt(1)
	for(d.Cmp(big.NewInt(1)) == 0) {
		x = f(x) 
		y = f(f(y))
		fmt.Println(x)
		fmt.Println(y)
		fmt.Println()
		result := new(big.Int)
		result.Sub(x,y)
		result.Abs(result)
		d = result.GCD(nil, nil, result, toFactor)
	}
	if(d.Cmp(toFactor) == 0) {
		return d, true
	}
	
	return d, false
}

func get_f(toFactor *big.Int) polynomial {
	return func(x *big.Int) *big.Int {
		result := new(big.Int).Mul(x,x)
		result.Sub(result, big.NewInt(1))
		result.Mod(result, toFactor)
		return result
	}
}

func factorise(toFactor *big.Int) {	
	if(toFactor.ProbablyPrime(20)) {
		fmt.Println(toFactor)
		return
	}
	
	f := get_f(toFactor) //big.NewInt(rng.Int63())
	factor,result := pollardRho(toFactor, f)
	
	fmt.Println(factor)
	fmt.Println(result)
}

func main() {
	toFactor := new(big.Int)
	//fmt.Println(rng.Int31())
	fmt.Scan(toFactor)		
	factorise(toFactor)
}
