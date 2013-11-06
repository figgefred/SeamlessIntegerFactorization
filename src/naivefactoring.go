
package  main

import "math/big"

var (
	primes = [...]*big.Int { big.NewInt(2), big.NewInt(3), big.NewInt(5), big.NewInt(7), big.NewInt(11), big.NewInt(13), big.NewInt(17), big.NewInt(19), big.NewInt(23), big.NewInt(29), big.NewInt(37), big.NewInt(41), big.NewInt(43), big.NewInt(47), big.NewInt(53), big.NewInt(59), big.NewInt(61), big.NewInt(67), big.NewInt(71), big.NewInt(73), big.NewInt(79), big.NewInt(83), big.NewInt(89), big.NewInt(97), big.NewInt(101)}
	// primes = [...]*big.Int { big.NewInt(2), big.NewInt(3), big.NewInt(5), big.NewInt(7), big.NewInt(11), big.NewInt(13), big.NewInt(17), big.NewInt(19), big.NewInt(23), big.NewInt(29)}
)

func trialdivision(task *Task, toFactor* big.Int) ([]*big.Int, *big.Int, bool) {
	factor := new(big.Int).Set(toFactor)
	resultBuffer := make([]*big.Int, 0, len(primes)+1)
	if factor.ProbablyPrime(prime_precision) {
		return append(resultBuffer, factor), nil, false
	}
	ZERO := big.NewInt(0)
	r := new(big.Int)

	hasDivided := true
	for hasDivided {
		hasDivided = false
		if(task.ShouldStop()) {
			return resultBuffer, factor, true
		}
		
		for _, prime := range primes {
			newFactor := new(big.Int)			
			newFactor.QuoRem(factor, prime, r)
			if r.Cmp(ZERO) != 0 {
				continue
			}

			resultBuffer = append(resultBuffer, prime)
			if newFactor.ProbablyPrime(prime_precision) {
				resultBuffer = append(resultBuffer, newFactor)
				factor = nil
				break
			} else {
				factor = newFactor
				hasDivided = true
				break
			}
		}
	}
	return resultBuffer, factor, false
}
