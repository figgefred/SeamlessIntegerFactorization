package main

import "fmt"
import "math/big"
import "math/rand"
import "time"
import "runtime"
import "bufio"
import "os"
import "strings"

type polynomial func(*big.Int) *big.Int

var (
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	deadline = 5
	prime_precision = 20
)

func pollardRho(toFactor *big.Int, f polynomial) (*big.Int, bool) {
	var x,y,d *big.Int
	x = big.NewInt(2)
	y = big.NewInt(2)
	d = big.NewInt(1)
	for(d.Cmp(big.NewInt(1)) == 0) {
		x = f(x) 
		y = f(f(y))
		//~ fmt.Println(x)
		//~ fmt.Println(y)
		//~ fmt.Println()
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
		result.Add(result, big.NewInt(rng.Int63()))
		result.Mod(result, toFactor)
		return result
	}
}

func factorise(toFactor *big.Int) string {	
	if(toFactor.ProbablyPrime(prime_precision)) {
		return toFactor.String() + "\n"
	}
	
	quo := new(big.Int)
	quo.Set(toFactor)
	
	var buf string
	//f := get_f(new(big.Int).Mul(toFactor,toFactor))
	for(quo.Cmp(big.NewInt(1)) > 0) {
		f := get_f(toFactor)
		factor,error := pollardRho(quo, f)
		
		if(error) {
			// Try again
			continue
		}
		
		quo.Quo(quo, factor)				
		
		if(!factor.ProbablyPrime(prime_precision)) {
			buf += factorise(factor)
		} else {
			buf += factor.String() + "\n" 
		}
		
		if(quo.ProbablyPrime(prime_precision)) {
			buf += quo.String() + "\n"
			break
		}
	}
	
	return buf
}

var result chan string

func main() {
		
	reader := bufio.NewReader(os.Stdin)

	//factorCount := 100
	factorCount := 2
	factorValues := make([]*big.Int, factorCount)
	// Read in line by line
    for i := 0; i < factorCount; i++ {
        line, _ := reader.ReadString('\n')
        factorValues[i] = new(big.Int)
        if _, ok := factorValues[i].SetString(strings.TrimSpace(line), 10) ; !ok {
        	fmt.Println("Parse error of", line)
        	
			// Exit
        	return
        }
    }

	runtime.GOMAXPROCS(runtime.NumCPU())
	
	for _, toFactor := range factorValues {

		//timeout := time.After(time.Duration(deadline) * time.Second)
		result = make(chan string)
		timeout := make(chan bool, 1)
		
		go func() {
			result <- factorise(toFactor)
		}();
		
		go func() {
			time.Sleep(time.Duration(deadline) * time.Second)	
			timeout <- true
		}();
		
		select {		
			case factors := <- result:
				fmt.Println(factors)	
			case <- timeout:
				fmt.Println("fail")	
				fmt.Println()		
		}
	}	
}
