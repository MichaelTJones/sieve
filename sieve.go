// Package sieve implements the prime sieve of Eratosthenes with clarity.
// Optimizations include odd-number tally, bit packing, and large factor
// termination. Prime sieves up to 1,000,000,000 are built quickly. A few
// prime-related functions are also provided. (This is not a segmented
// wheel implementation)
package sieve

import (
	"fmt"
	"math"
)

type word uint8

const wordBitsLog2 = 3

// type word uint8 ; const wordBitsLog2 = 3
// type word uint16 ; const wordBitsLog2 = 4
// type word uint32 ; const wordBitsLog2 = 5
// type word uint64 ; const wordBitsLog2 = 6

const wordBits = 1 << wordBitsLog2
const wordMask = wordBits - 1

/*
type word uint8

const (
	wordMask = uint(^word(0)) // needed
	wordBytesLog2 = wordMask>>8&1 + wordMask>>16&1 + wordMask>>32&1
	wordBytes = 1<<uint(wordBytesLog2)
	wordBitsLog2 = wordBytesLog2 + 2 // needed
	wordBits = 1<<uint(wordBitsLog2) // needed
)
*/

type Sieve struct {
	size  int    // the largest number testable for primality
	count int    // the number of primes resident in the sieve
	table []word // the sieve, one bit per odd number >= 3
}

// bit gets the value of bit[index] by inspecting bits in a packed table.
func (sieve *Sieve) bit(index int) (bit byte) {
	return byte(sieve.table[index>>(1+wordBitsLog2)]>>byte((index>>1)&wordMask)) & 0x1 // bit storage
}

// setBit sets the value of bit[index] to one (marks it as non-prime) in a packed table.
func (sieve *Sieve) setBit(index int) {
	sieve.table[index>>(1+wordBitsLog2)] |= word(1 << ((uint(index) >> 1) & wordMask)) // bit storage
	return
}

// New allocates and initializes a prime sieve representing all primes <= size using
// the method of Eratosthenes of Cyrene (Libya.) This approach has prevailed for 2207
// years and vies with Euclid of Alexandria's GCD method as the first known algorithm.
func New(size int) *Sieve {
	sieve := new(Sieve)
	sieve.size = size
	sieve.table = make([]word, 1+(sieve.size+wordBits-1)/wordBits) // one bit per odd number in range
	for i := 3; i*i <= sieve.size; i += 2 {                        // early exit for the larger factor
		if sieve.bit(i) == 0 { // next prime
			for j := 3 * i; j <= sieve.size; j += i + i {
				sieve.setBit(j) // strike multiples from table
			}
		}
	}
	return sieve
}

// NewCount allocates and initializes a sieve sized to include the first count primes.
func NewCount(count int) *Sieve {
	// estimate size from the prime number theorem's asymptotic value
	size := 1.25506 * float64(count) * math.Log(float64(count))
	if size < 64 {
		size = 64 // increase estimate for small values of count
	}
	return New(int(size))
}

// NewFactor allocates and initializes a sive sized to factor numbers <= n.
func NewFactor(n int) *Sieve {
	size := math.Ceil(math.Sqrt(float64(n)))
	return New(int(size) + 32) // extra 32 is optional cushion
}

// Size returns the largest number whose primality is encoded in the sieve. This
// number need not be a prime. Used to range-check the sieve before prime testing.
// The companion Factor() method handles numbers up to Size()*Size().
func (sieve *Sieve) Size() int {
	return sieve.size
}

// Count the number of primes in the sieve.
func (sieve *Sieve) Count() int {
	if sieve.count == 0 {
		if sieve.size >= 2 {
			sieve.count = 1 // account for first prime, 2, not in table
		}
		for i := 3; i <= sieve.size; i += 2 {
			if sieve.bit(i) == 0 { // is prime
				sieve.count++
			}
		}
	}
	return sieve.count
}

// Prime tests primality using the sieve for precomputed answer. Testing values outside
// the range of the sieve returns false, indicating that this sieve cannot prove the
// number to be prime. For values inside the range, the result is definitive.
func (sieve *Sieve) Prime(n int) bool {
	switch {
	case n < 2:
		return false
	case n == 2:
		return true
	case n&1 == 0:
		return false
	case n <= sieve.size:
		// determine primality by direct inspection
		return n == 2 || (n > 2 && n <= sieve.size && n&1 == 1 && sieve.bit(n) == 0)
	case n <= sieve.size*sieve.size:
		// determine primality by trial division
		root := 1
		for (root+1)*(root+1) <= n {
			root++
		}
		if root*root == n { // perfect square
			return false
		}
		for d := 3; d <= root; d += 2 {
			if sieve.bit(d) == 0 && n%d == 0 {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// String returns a string of the sieve's primes useful to functions in the fmt package.
// It matches the Stringer interface to support output with fmt's "%v" and "%s" modes or
// directly, as in p := sieve.New(100); fmt.Println(p)
func (sieve *Sieve) String() string {
	var s string
	first := true
	if sieve.size >= 2 {
		s += "2"
		first = false
	}
	for i := 3; i <= sieve.size; i += 2 {
		if sieve.bit(i) == 0 { // next prime
			if !first {
				s += " "
			}
			s += fmt.Sprintf("%d", i)
			first = false
		}
	}
	return s
}

// Factor an integer <= sieve.Size()*sieve.Size() using the sieve for trial divisors.
// Returns a slice of factors. Repeated factors are repeated in the result.
func (sieve *Sieve) Factor(n int) []int {
	if sieve.size*sieve.size < n { // too big for sieve?
		return make([]int, 0, 0)
	}
	if n <= 3 {
		result := make([]int, 1, 1)
		result[0] = n
		return result
	}

	result := make([]int, 0, 64)
	for n > 1 && n%2 == 0 { // initial factor of 2
		result = append(result, 2)
		// n /= 2
		n >>= 1
	}
	for d := 3; d < sieve.size && d*d <= n; d += 2 {
		if sieve.Prime(d) { // try primes in 3..sqrt(n)
			for n > 1 && n%d == 0 {
				result = append(result, d)
				n /= d
			}
		}
	}
	if n > 1 { // remaining prime factor
		result = append(result, n)
		n = 1
	}
	return result
}

type Unique struct {
	Factor int
	Count  int
}

// Factor an integer <= sieve.Size()*sieve.Size() using the sieve for trial divisors.
// Returns a slice of factors. Repeated factors are repeated in the result.
func (sieve *Sieve) FactorUnique(n int) []Unique {
	if sieve.size*sieve.size < n { // too big for sieve?
		return make([]Unique, 0, 0)
	}
	if n <= 3 {
		result := make([]Unique, 1, 1)
		result[0].Factor = n
		result[0].Count = 1
		return result
	}

	result := make([]Unique, 0, 64)
	if n > 1 && n%2 == 0 { // initial factor of 2
		count := 0
		for n > 1 && n%2 == 0 {
			// n /= 2
			n >>= 1
			count++
		}
		result = append(result, Unique{2, count})
	}
	for d := 3; d < sieve.size && d*d <= n; d += 2 {
		if sieve.Prime(d) { // try primes in 3..sqrt(n)
			if n > 1 && n%d == 0 {
				count := 0
				for n > 1 && n%d == 0 {
					n /= d
					count++
				}
				result = append(result, Unique{d, count})
			}
		}
	}
	if n > 1 { // remaining prime factor
		result = append(result, Unique{n, 1})
		n = 1
	}
	return result
}

func (sieve *Sieve) FactorString(n int) string {
	u := sieve.FactorUnique(n)
	r := ""
	space := ""
	for _, f := range u {
		if f.Count == 1 {
			r = r + fmt.Sprintf("%s%d", space, f.Factor)
		} else {
			r = r + fmt.Sprintf("%s%d^%d", space, f.Factor, f.Count)
		}
		space = " "
	}
	return r
}

// Determine the total number of divisors of n
// Divisors(6) == 4, from {1, 2, 3, 6}
func (sieve *Sieve) DivisorCount(n int) int {
	if sieve.size*sieve.size < n { // too big for sieve?
		return 0
	}
	if n == 1 {
		return 1
	}
	// if n <= 3 {
	// 	return 2
	// }
	m := 1
	if n > 1 && n%2 == 0 { // initial factors of 2
		count := 0
		for n > 1 && n%2 == 0 {
			// n /= 2
			n >>= 1
			count++
		}
		m *= count + 1
	}
	for d := 3; d < sieve.size && d*d <= n; d += 2 {
		if sieve.Prime(d) { // try primes in 3..sqrt(n)
			if n > 1 && n%d == 0 {
				count := 0
				for n > 1 && n%d == 0 {
					n /= d
					count++
				}
				m *= count + 1
			}
		}
	}
	if n > 1 { // remaining prime factor
		m *= (1 + 1)
	}
	return m
}

// SquareFree is a boolean test that the subject number's factors are not repeated.
// Square-free numbers are the sequence http://oeis.org/A005117
func (sieve *Sieve) SquareFree(n int) bool {
	if sieve.size*sieve.size < n { // too big for sieve?
		return false
	}
	f := sieve.Factor(n)
	for i := 0; i < len(f)-1; i++ {
		if f[i] == f[i+1] { // is there a repeated factor?
			return false
		}
	}
	return true
}

func (sieve *Sieve) NthPrime(n int) int {
	if n == 1 {
		return 2
	}
	for value, count := 3, 1; value <= sieve.size; value += 2 {
		if sieve.Prime(value) {
			count++
			if count == n {
				return value
			}
		}
	}
	return 0 // this sieve contains less than n primes
}
