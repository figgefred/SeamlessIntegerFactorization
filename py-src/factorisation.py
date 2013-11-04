#!/usr/bin/python
import math
import numbers
import fileinput
import sys

# Miller-Rabin primality test
# http://en.wikipedia.org/wiki/Miller%E2%80%93Rabin_primality_test
def isPrime(number):	
	# If the number is even and not 2 it is not a prime.
	if number > 2 and (number & 1) < 1: 
		return False	
	
	# 2 is a prime (special case)	
	alist = {2,3,5,7,11,13,17,19,23}	
	if number in alist:
		return True
		
	# write n-1 as 2s*d by factoring powers of 2 from n-1	
	s = 0
	d = number-1
	while (long(math.ceil(d)) & 1) < 1:
		s += 1
		d /= 2
	
	
	for a in alist:
		# Fermats theorem, if a^d = 1 mod number then number is coprime!		
		if pow(a,long(math.ceil(d)),number) != 1:		
			truedat = True			
			for r in range(s):
				truedat &= (pow(a, 2**r * long(math.ceil(d)), number) != number-1)								
			if(truedat):				
				return False			
	return True
	
def isCoprime(number1, number2):
	return number1 % number2 > 0

primes = set({2,3,5,7,11,13,17,19})
maxprime = 17
primes.add(maxprime)

while 1:
	try:
		line = sys.stdin.readline()
		
		if not line or len(line) <= 1:
			break
			
	except KeyboardInterrupt:
		print "interuppted"
		break

	target = long(line)
	if isPrime(target):
		print target		
		print
		primes.add(target)
		continue

	print "fail"
	print 
	"""
	factors = []	
	for factorMaybe in primes:
		while not isCoprime(target, factorMaybe) and target > 1:		
			target = target / factorMaybe
			factors.append(factorMaybe)
			if isPrime(target):
				factors.append(target)
				target = 1
				break
		

	factorMaybe = maxprime + 2	
	while target > 1 and factorMaybe < target / maxprime and len(primes) < 50000:		
		if isPrime(factorMaybe):
			maxprime = factorMaybe
			primes.add(maxprime)
			while not isCoprime(target, factorMaybe) and target > 1:		
				target = target / factorMaybe
				factors.append(factorMaybe)
				if isPrime(target):
					factors.append(target)
					target = 1
					break				
			
			
		factorMaybe += 2

	if target == 1:
		for factor in factors:
			print factor
	else:
		print "fail"	
	print
"""	
	
		
	
		
	
		
		
