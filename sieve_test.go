package sieve

import (
	"fmt"
	"testing"
)

//
// TESTS
//

// Basic sanity test: are primes up to 100 determined properly?
func TestSanity(t *testing.T) {
	result := New(100).String()
	expect := "2 3 5 7 11 13 17 19 23 29 31 37 41 43 47 53 59 61 67 71 73 79 83 89 97"
	if result != expect {
		t.Errorf("New(100).String() is %q; want %q", result, expect)
	}
}

var countTests = []struct {
	size  int
	count int
}{
	{10, 4},          // π(10^1) antiquity
	{100, 25},        // π(10^2) L. Pisano (1202; Beiler)
	{1000, 168},      // π(10^3) F. van Schooten (1657; Beiler)
	{10000, 1229},    // π(10^4) F. van Schooten (1657; Beiler)
	{100000, 9592},   // π(10^5) T. Brancker (1668; Beiler)
	{1000000, 78498}, // π(10^6) A. Felkel (1785; Beiler)
	// {10000000, 664579},   // π(10^7) J. P. Kulik (1867; Beiler)
	// {100000000, 5761455}, // π(10^8) Meissel (1871; corrected)
	// {153339973, 8621475},
	// {1000000000, 50847534},       // π(10^9) Meissel (1886; corrected)
	// {10000000000, 455052511},     // π(10^10)
	// {100000000000, 4118054813},   // π(10^11)
	// {1000000000000, 37607912018}, // π(10^12)
}

// Are the number of sieve-surviving primes <= n equal to π(n) as expected?
func TestCounts(t *testing.T) {
	for i, a := range countTests {
		count := New(a.size).Count()
		if count != a.count {
			t.Errorf("#%d, New(%d).Count() is %d; want %d", i, a.size, count, a.count)
		}
	}
}

var sumTests = []struct {
	count int    // number of primes to sum
	index int    // index of the count'th prime (http://primes.utm.edu/nthprime)
	sum   uint64 // sum of the first count primes (http://oeis.org/A007504)
}{
	{10, 29, 129},
	{100, 541, 24133},
	{1000, 7919, 3682913},
	{10000, 104729, 496165411},
	{100000, 1299709, 62260698721},
	{1000000, 15485863, 7472966967499},
	// {10000000, 179424673, 870530414842019},
	// {100000000, 2038074743, 99262851056183695},
}

// Sum first n primes (or return 0 if sieve contains less than n primes)
func (sieve *Sieve) sum(n int) (sum uint64) {
	if n >= 2 { // only even prime, 2, is not in table
		sum += 2
		n--
	}
	for i := 3; i <= sieve.size && n > 0; i += 2 {
		if sieve.bit(i) == 0 {
			sum += uint64(i)
			n--
		}
	}
	// was sieve too small?
	if n > 0 {
		sum = 0 // indicate failure
	}
	return
}

// Are the sums of the first few primes consistent with expectation?
func TestSums(t *testing.T) {
	for i, a := range sumTests {
		sieve := New(a.index)
		sum := sieve.sum(a.count)
		if sum != a.sum {
			t.Errorf("#%d, sum of first %d primes = %d; want %d", i, a.count, sum, a.sum)
		}
	}
}

var twinTests = []struct {
	size  int // primes <= size
	twins int // number of twin primes (n and n+2 are both prime)
}{
	{10, 2},
	{100, 8},
	{1000, 35},
	{10000, 205},
	{100000, 1224},
	{1000000, 8169},
	// {10000000, 58980},
	// {100000000, 440312},
	// {1000000000, 3424506},
	// {10000000000, 27412679},
	// {100000000000, 224376048},
	// {1000000000000, 1870585220},
	// {10000000000000, 15834664872},
	// {100000000000000, 135780321665},
	// {1000000000000000, 1177209242304},
	// {10000000000000000, 19831847025792}, // see Thomas R. Nicely, http://www.trnicely.net
}

// Are the computed number of twin primes in the range 3 .. a.size equal to the expected result?
func TestTwins(t *testing.T) {
	for i, a := range twinTests {
		s := New(a.size + 2)
		twins := 0
		for i := 3; i <= a.size; i += 2 {
			if s.Prime(i) == true && s.Prime(i+2) == true {
				twins++
			}
		}
		if twins != a.twins {
			t.Errorf("#%d, number of twin primes in first %d primes = %d; want %d", i, a.size, twins, a.twins)
		}
	}
}

var triple024Tests = []struct {
	size    int // primes <= size
	triples int // number of "024" triple primes (n+0 and n+2 and n+4 are all prime)
}{
	{10, 1},
	{100, 1},
	{1000, 1},
	{10000, 1},
	{100000, 1},
	{1000000, 1},
	// {10000000, 1},
	// {100000000, 1},
	// {1000000000, 1},
	// {10000000000, 1},
	// {100000000000, 1},
	// {1000000000000, 1},
	// {10000000000000, 1},
	// {100000000000000, 1},
	// {1000000000000000, 1},
	// {10000000000000000, 1},
}

// Are the computed number of "024" triple primes in the range 3 .. a.size equal to the expected result?
// [3,5,7] is the only solution
func Test024Triples(t *testing.T) {
	for i, a := range triple024Tests {
		s := New(a.size + 4)
		triples := 0
		for i := 3; i <= a.size; i += 2 {
			if s.Prime(i) == true && s.Prime(i+2) == true && s.Prime(i+4) == true {
				triples++
			}
		}
		if triples != a.triples {
			t.Errorf("#%d, number of 024 triple primes in first %d primes = %d; want %d", i, a.size, triples, a.triples)
		}
	}
}

var triple026Tests = []struct {
	size    int // primes <= size
	triples int // number of "026" triple primes (n+0 and n+2 and n+6 are all prime)
}{
	{10, 1},
	{100, 4},
	{1000, 15},
	{10000, 55},
	{100000, 259},
	{1000000, 1393},
	// {10000000, 8543},
	// {100000000, 55600},
	// {1000000000, 379508},
	// {10000000000, 2713347},
	// {100000000000, 20093124},
	// {1000000000000, 152850135},
	// {10000000000000, 1189795268},
	// {100000000000000, 9443899421},
	// {1000000000000000, 76218094021},
	// {10000000000000000, 624026299748}, // see Thomas R. Nicely, http://www.trnicely.net
}

// Are the computed number of "026" triple primes in the range 3 .. a.size equal to the expected result?
func Test026Triples(t *testing.T) {
	for i, a := range triple026Tests {
		s := New(a.size + 6)
		triples := 0
		for i := 3; i <= a.size; i += 2 {
			if s.Prime(i) == true && s.Prime(i+2) == true && s.Prime(i+6) == true {
				triples++
			}
		}
		if triples != a.triples {
			t.Errorf("#%d, number of 026 triple primes in first %d primes = %d; want %d", i, a.size, triples, a.triples)
		}
	}
}

var triple046Tests = []struct {
	size    int // primes <= size
	triples int // number of "046" triple primes (n+0 and n+4 and n+6 are all prime)
}{
	{10, 1},
	{100, 5},
	{1000, 15},
	{10000, 57},
	{100000, 248},
	{1000000, 1444},
	// {10000000, 8677},
	// {100000000, 55556},
	// {1000000000, 379748},
	// {10000000000, 2712226},
	// {100000000000, 20081601},
	// {1000000000000, 152839134},
	// {10000000000000, 1189826966},
	// {100000000000000, 9443942237},
	// {1000000000000000, 76217933571},
	// {10000000000000000, 624025508307}, // see Thomas R. Nicely, http://www.trnicely.net
}

// Are the computed number of "046" triple primes in the range 3 .. a.size equal to the expected result?
func Test046Triples(t *testing.T) {
	for i, a := range triple046Tests {
		s := New(a.size + 6)
		triples := 0
		for i := 3; i <= a.size; i += 2 {
			if s.Prime(i) == true && s.Prime(i+4) == true && s.Prime(i+6) == true {
				triples++
			}
		}
		if triples != a.triples {
			t.Errorf("#%d, number of 046 triple primes in first %d primes = %d; want %d", i, a.size, triples, a.triples)
		}
	}
}

// Count the number of primes in the form 2*(n**2)-1 for 1 <= n <= upper
// This is Project Euler problem 216: http://projecteuler.net/problem=216
// though not solved the way optimal for the upper = 50,000,000 case

const sqrt2 = 1.4142135623730950488016887242096980785696718753769480731766797379907324784621

func countEuler(limit int) int {
	size := 1 + int(sqrt2*float64(limit)) // minimal sieve for factors
	//size := 2*limit // nearly minimal sieve for factors
	//size := 2*limit*limit-1 // complete sieve for candidates
	s := New(size) // build sieve
	count := 0
	for n := 2; n <= limit; n++ {
		if s.Prime(2*n*n-1) == true {
			count++
		}
	}
	return count
}

var eulerTests = []struct {
	upper int // test 2*n*n+1 for 2 <= n <= upper
	count int // number of polynomial values that are prime
}{
	{10, 7},
	{100, 45},
	{1000, 303},
	{10000, 2202},
	// {100000, 17185},
	// {1000000, 141444},
}

func TestEuler(t *testing.T) {
	for i, a := range eulerTests {
		count := countEuler(a.upper)
		if count != a.count {
			t.Errorf("#%d, number of primes 2*(n**2)-1, n in 1..%d = %d; want %d", i, a.upper, count, a.count)
		}
	}
}

// Test Nth Primes
var nth = []struct {
	n     int
	prime int
}{
	{10, 29},
	{100, 541},
	{1000, 7919},
	{10000, 104729},
	{100000, 1299709},
	// {1000000, 15485863},
	// {8621475, 153339973},
}

func TestNthPrime(t *testing.T) {
	for i, a := range nth {
		s := NewCount(a.n)
		p := s.NthPrime(a.n)
		if p != a.prime {
			t.Errorf("#%d, prime[%v] = %d; want %d", i, a.n, p, a.prime)
		}
	}
}

//
// BENCHMARKS
//

// Measure the average time it takes to generate sieves of various sizes.
// This grows as O(n log n) in Eratosthenes' algorithm, or slightly slower
// thanks to the "i*i <= sieve.size" early termination in the outer loop.
// In practice, the actual measured time is dominted by cache hieracrhy in
// memory access. Also, the getter and setter functions cost 5%-8% more
// than doing it inline, but being inline obfuscated the logic.

func benchmarkNew(b *testing.B, n int) {
	for i := 0; i < b.N; i++ {
		_ = New(n)
	}
}

// func BenchmarkNew153339973(b *testing.B) { benchmarkNew(b, 153339973) }

func BenchmarkNew10(b *testing.B)      { benchmarkNew(b, 10) }
func BenchmarkNew100(b *testing.B)     { benchmarkNew(b, 100) }
func BenchmarkNew1000(b *testing.B)    { benchmarkNew(b, 1000) }
func BenchmarkNew10000(b *testing.B)   { benchmarkNew(b, 10000) }
func BenchmarkNew65536(b *testing.B)   { benchmarkNew(b, 65536) }
func BenchmarkNew100000(b *testing.B)  { benchmarkNew(b, 100000) }
func BenchmarkNew1000000(b *testing.B) { benchmarkNew(b, 1000000) }

// func BenchmarkNew10000000(b *testing.B)      { benchmarkNew(b, 10000000) }
// func BenchmarkNew100000000(b *testing.B)     { benchmarkNew(b, 100000000) }
// func BenchmarkNew1000000000(b *testing.B)    { benchmarkNew(b, 1000000000) }
// func BenchmarkNew10000000000(b *testing.B)   { benchmarkNew(b, 10000000000) }
// func BenchmarkNew100000000000(b *testing.B)  { benchmarkNew(b, 100000000000) }
// func BenchmarkNew1000000000000(b *testing.B) { benchmarkNew(b, 1000000000000) }

// Measure the average time it takes to perform a primality test in sieves
// of various sizes. The sieve structure makes this fast and O(1), at a
// storage cost of Size()/16 bytes to store the encoded sieve.

func benchmarkPrime(b *testing.B, n int) {
	b.StopTimer()
	s := New(n)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = s.Prime(n)
	}
}

func BenchmarkPrime10(b *testing.B)      { benchmarkPrime(b, 10) }
func BenchmarkPrime100(b *testing.B)     { benchmarkPrime(b, 100) }
func BenchmarkPrime1000(b *testing.B)    { benchmarkPrime(b, 1000) }
func BenchmarkPrime10000(b *testing.B)   { benchmarkPrime(b, 10000) }
func BenchmarkPrime65536(b *testing.B)   { benchmarkPrime(b, 65536) }
func BenchmarkPrime100000(b *testing.B)  { benchmarkPrime(b, 100000) }
func BenchmarkPrime1000000(b *testing.B) { benchmarkPrime(b, 1000000) }

// func BenchmarkPrime10000000(b *testing.B)   { benchmarkPrime(b, 10000000) }
// func BenchmarkPrime100000000(b *testing.B)  { benchmarkPrime(b, 100000000) }
// func BenchmarkPrime1000000000(b *testing.B) { benchmarkPrime(b, 1000000000) }

// Measure the average time to factor integers between 1 and 1,000,000
func BenchmarkFactor(b *testing.B) {
	rangeLow := 2
	rootHigh := 1000
	rangeHigh := rootHigh * rootHigh
	delta := rangeHigh - rangeLow

	b.StopTimer()
	s := New(rootHigh)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = s.Factor(rangeLow + i%delta)
	}
}

//
// EXAMPLES
//

func ExampleSieve() {
	// Build sieve for primes up to 10.
	s := New(10)
	fmt.Println(s)
	// Output:
	// 2 3 5 7
}

func ExampleSieve_Count() {
	// Count number of primes up to 10.
	s := New(10)
	fmt.Println(s.Count())
	// Output:
	// 4
}

func ExampleSieve_Prime() {
	// Test numbers for primality.
	m := New(1000000)
	fmt.Println(8, m.Prime(8))
	fmt.Println(999983, m.Prime(999983))
	// Output:
	// 8 false
	// 999983 true
}

func ExampleSieve_Factor() {
	// Factor numbers up to 10
	s := New(10)
	for i := 2; i <= 10; i++ {
		fmt.Println(i, s.Factor(i))
	}
	// Output:
	// 2 [2]
	// 3 [3]
	// 4 [2 2]
	// 5 [5]
	// 6 [2 3]
	// 7 [7]
	// 8 [2 2 2]
	// 9 [3 3]
	// 10 [2 5]
}

func ExampleSieve_FactorUnique() {
	// Factor numbers up to 10
	s := New(10)
	for i := 2; i <= 10; i++ {
		fmt.Println(i, s.FactorUnique(i))
	}
	// Output:
	// 2 [{2 1}]
	// 3 [{3 1}]
	// 4 [{2 2}]
	// 5 [{5 1}]
	// 6 [{2 1} {3 1}]
	// 7 [{7 1}]
	// 8 [{2 3}]
	// 9 [{3 2}]
	// 10 [{2 1} {5 1}]
}

func ExampleSieve_SquareFree() {
	// Determine square-free numbers up to 10
	s := New(10)
	for i := 2; i <= 10; i++ {
		fmt.Println(i, s.SquareFree(i))
	}
	// Output:
	// 2 true
	// 3 true
	// 4 false
	// 5 true
	// 6 true
	// 7 true
	// 8 false
	// 9 false
	// 10 true
}
