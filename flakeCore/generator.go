package flakeCore

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

// prime.go - Generate safe primes
//
// Copyright 2013-2017 Sudhi Herle <sudhi.herle-at-gmail-dot-com>
// License: MIT

var one *big.Int

var simplePrimes = []int64{
	2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61,
	67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109, 113, 127, 131, 137,
	139, 149, 151, 157, 163, 167, 173, 179, 181, 191, 193, 197, 199, 211,
	223, 227, 229, 233, 239, 241, 251, 257, 263, 269, 271, 277, 281, 283,
	293, 307, 311, 313, 317, 331, 337, 347, 349, 353, 359, 367, 373, 379,
	383, 389, 397, 401, 409, 419, 421, 431, 433, 439, 443, 449, 457, 461,
	463, 467, 479, 487, 491, 499, 503, 509, 521, 523, 541,
}

var bitRange = []int{
	128, 256, 512,
}

type PrimeField struct {
	BG     *big.Int
	BN     *big.Int
	SN     int // Length of N in bytes
	Length int
}

//Generates N and G
func newPrimeField(nbits int) (*PrimeField, error) {

	for i := 0; i < 100; i++ {
		p, err := safePrime(nbits)
		if err != nil {
			return nil, err
		}

		for _, g0 := range simplePrimes {
			g := big.NewInt(g0)
			if isGenerator(g, p) {
				pf := &PrimeField{
					BG:     g,
					BN:     p,
					SN:     nbits / 8,
					Length: nbits,
				}
				return pf, nil
			}
		}
	}
	return nil, fmt.Errorf("srp: can't find generator after 100 tries")
}

func safePrime(bits int) (*big.Int, error) {

	a := new(big.Int)
	for {
		p, err := rand.Prime(rand.Reader, bits)
		if err != nil {
			return nil, err
		}

		// 2p+1
		a = a.Lsh(p, 1)
		a = a.Add(a, one)
		if a.ProbablyPrime(20) {
			return a, nil
		}
	}
	return nil, nil
}

func isGenerator(g, p *big.Int) bool {
	p1 := big.NewInt(0).Sub(p, one)
	q := big.NewInt(0).Rsh(p1, 1) // q = p-1/2 = ((p-1) >> 1)

	// p is a safe prime. i.e., it is of the form 2q+1 where q is prime.
	//
	// => p-1 = 2q, where q is a prime.
	//
	// All factors of p-1 are: {2, q, 2q}
	//
	// So, our check really comes down to:
	//   1) g ^ (p-1/2q) != 1 mod p
	//		=> g ^ (2q/2q) != 1 mod p
	//		=> g != 1 mod p
	//	    Trivial case. We ignore this.
	//
	//   2) g ^ (p-1/2) != 1 mod p
	//      => g ^ (2q/2) != 1 mod p
	//      => g ^ q != 1 mod p
	//
	//   3) g ^ (p-1/q) != 1 mod p
	//      => g ^ (2q/q) != 1 mod p
	//      => g ^ 2 != 1 mod p
	//

	// g ^ 2 mod p
	if !ok(g, big.NewInt(0).Lsh(one, 1), p) {
		return false
	}

	// g ^ q mod p
	if !ok(g, q, p) {
		return false
	}

	return true
}

func ok(g, x *big.Int, p *big.Int) bool {
	z := big.NewInt(0).Exp(g, x, p)
	// the expmod should NOT be 1
	return z.Cmp(one) != 0
}

func hex2str(n *big.Int) string {
	return fmt.Sprintf("%X", n) // or %X or upper case
}

func str2hex(s string) *big.Int {
	x := new(big.Int)
	x.SetString(s[0:len(s)], 16)
	return x
}

func (p *PrimeField) Converstring() string {
	return strconv.Itoa(p.Length) + ":" + strconv.Itoa(p.SN) + ":" + hex2str(p.BG) + ":" + hex2str(p.BN)
}

func (p *PrimeField) GenK() (*big.Int, string) {
	strA := hex2str(p.BN)
	strB := hex2str(p.BG)
	sbyte := []byte(strA + ":" + strB)
	sHash := sha256.Sum256(sbyte)
	hexStr := hex.EncodeToString(sHash[:])
	bint := str2hex(hexStr)
	hexStr = strings.ToUpper(hexStr)
	return bint, hexStr
}

//func writeFile( x string  , path string ) error {
//	os.
//}

func Generate() []*PrimeField {
	arr := []*PrimeField{}
	for _, s := range bitRange {
		a, _ := newPrimeField(s)
		arr = append(arr, a)
	}
	return arr
}

func A2S(pfRange []*PrimeField) []*string {
	str := []*string{}
	for _, s := range pfRange {
		x := s.Converstring()
		str = append(str, &x)
	}
	return str
}

func S2A(strRange []*string) []*PrimeField {
	pfRange := []*PrimeField{}
	for _, line := range strRange {
		v := strings.Split(*line, ":")
		n, _ := strconv.Atoi(v[1])
		s, _ := strconv.Atoi(v[0])
		pf := &PrimeField{
			BG:     str2hex(v[2]),
			BN:     str2hex(v[3]),
			SN:     n,
			Length: s,
		}
		pfRange = append(pfRange, pf)
	}
	return pfRange
}

func init() {
	one = big.NewInt(1)
	PrimeArray = S2A(pfStr)
	pc = PrimeArray
	ps = PrimeArray
}

func Str2PF(str string) *PrimeField {
	v := strings.Split(str, ":")
	s, _ := strconv.Atoi(v[0])
	n, _ := strconv.Atoi(v[1])
	pf := &PrimeField{
		Length: s, // s
		SN:     n, // n
		BG:     str2hex(v[2]),
		BN:     str2hex(v[3]),
	}
	return pf
}
