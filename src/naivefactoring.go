
package  main

import "math/big"

var (
	primes []uint16 = getPrimes() // Prime values of up to 65535 ... Count is 6544
	primeCalcCount = 5000		  // Above 6544 equates to 6544
								  // 2250 verkara vara Kattis mogen
)

func getPrimes() ([]uint16) {
	primes := make([]uint16, 0, primeCalcCount)
	//TWO := big.NewInt(2)
	primes = append(primes,2)
	val := 3
	primes = append(primes, uint16(val))
	
	for i := 2; i < primeCalcCount && val < 65535; i++ {
		for {
			if big.NewInt(int64(val)).ProbablyPrime(prime_precision) {
				primes = append(primes, uint16(val))
				break
			}
			val += 2
		}
		val += 2
	}
	//fmt.Println("FOUND", len(primes), "primes")
	return primes
}

/*func getPrimes() ([]*big.Int) {
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
}*/

// Returns true iff prime is a divisor of 'toFactor'
// Else false
// *big.Int will refer to an Int, yet is only guaranteed to be
// the true quotient if bool is true.
func divide(toFactor, prime *big.Int) (bool, *big.Int) {
		newFactor := new(big.Int)
		r := new(big.Int)
		newFactor.QuoRem(toFactor, prime, r)
		return r.Cmp(ZERO) == 0, newFactor
}

func naiveFactoring(task *Task) ([]*big.Int) {

	res, factor := trialdivision(task)
	if factor != nil {
		task.timed_out = true
	}
	return res
}

func trialDivisionPollardFactoring(task *Task) ([]*big.Int) {

	res, factor := trialdivision(task)
	if factor == nil {
		return res
	}
	res = append(res, _pollardFactoring(task, factor)...)
	return res
}

func trialdivision(task *Task) ([]*big.Int, *big.Int) {
	factor := task.toFactor
	resultBuffer := make([]*big.Int, 0, 20)
	if factor.ProbablyPrime(prime_precision) {
		return append(resultBuffer, factor), nil
	}
	// Loop over to find primes that divide 'factor'
	for _, p := range primes {
		// Prime greater than 'factor', then just break
		prime := big.NewInt(int64(p))
		if(task.ShouldStop()) {
			return resultBuffer, factor
		}

		if prime.Cmp(factor) > 0 {
			break
		}
		divisible, newFactor := divide(factor, prime)
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
	for _, prime := range resultBuffer {
		// Prime greater than 'factor', then just break
		if prime.Cmp(factor) > 0 {
			break
		}
		divided, newFactor := divide(factor, prime)		
		for divided {
			tmp = append(tmp, prime)
			if newFactor.ProbablyPrime(prime_precision) {
				tmp = append(tmp, newFactor)
				factor = nil
				break
			} else {
				factor = newFactor	
			}
			divided, newFactor = divide(factor, prime)		
		}
		if factor == nil {
			break
		}
		if(task.ShouldStop()) {
			return resultBuffer, factor
		}
	}
	
	resultBuffer = append(resultBuffer, tmp...)
	return resultBuffer, factor
}

