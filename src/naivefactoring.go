
package  main

import "math/big"

var (
	primes []*big.Int = getPrimes()
	primeCalcCount = 50000
)

func getPrimes() ([]*big.Int) {
	primes := make([]*big.Int, primeCalcCount)
	
	TWO := big.NewInt(2)
	primes[0] = TWO
	val := big.NewInt(3)
	primes[1] = val
	
	for i := 2; i < primeCalcCount; i++ {
		val = new(big.Int).Add(val, TWO)
		for {
			if val.ProbablyPrime(prime_precision) {
				primes[i] = val	
				break
			}
			val = new(big.Int).Add(val, TWO)
		}
	}

//fmt.Println(primes)
	return primes
}

// Returns true iff prime is a divisor of 'toFactor'
// Else false
// *big.Int will refer to an Int, yet is only guaranteed to be
// the true quotient if bool is true.
func isDivisible(toFactor, prime *big.Int) (bool, *big.Int) {
		
		newFactor := new(big.Int)
		r := new(big.Int)
		newFactor.QuoRem(toFactor, prime, r)
		return r.Cmp(ZERO) == 0, newFactor
}

func naivefactoring(task *Task) ([]*big.Int) {
	res, _ := trialdivision(task)
	return res
}

func trialdivision(task *Task) ([]*big.Int, *big.Int) {
	factor := task.toFactor
	resultBuffer := make([]*big.Int, 0, 20)
	if factor.ProbablyPrime(prime_precision) {
		return append(resultBuffer, factor), nil
	}

	// Loop over to find primes that divide 'factor'
	for _, prime := range primes {
		// Prime greater than 'factor', then just break

		if(task.ShouldStop()) {
			return resultBuffer, factor
		}

		if prime.Cmp(factor) > 0 {
			break
		}
		divisible, newFactor := isDivisible(factor, prime)
		if !divisible {
			continue
		}
		resultBuffer = append(resultBuffer, prime)
		if newFactor.ProbablyPrime(prime_precision) {
			resultBuffer = append(resultBuffer, newFactor)
			factor = nil
			break
		} else {
			factor = newFactor	
		}
	}

	if factor == nil {
		return resultBuffer, factor
	}

	// Check if the 'factors' found divide the 'factor'
	
	// Temporary list of results
	tmp := make([]*big.Int, 0, cap(resultBuffer))
	hasDivided := true	
	for hasDivided {
		hasDivided= false
		for _, prime := range resultBuffer {
			// Prime greater than 'factor', then just break
			if prime.Cmp(factor) > 0 {
				break
			}
			divisible, newFactor := isDivisible(factor, prime)
			if !divisible {
				continue
			}
			tmp = append(tmp, prime)
			if newFactor.ProbablyPrime(prime_precision) {
				tmp = append(tmp, newFactor)
				factor = nil
				break
			} else {
				factor = newFactor	
			}
		}
		if(task.ShouldStop()) {
			return resultBuffer, factor
		}
	}
	resultBuffer = append(resultBuffer, tmp...)
	return resultBuffer, factor
}
