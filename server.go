package flakeCore

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

var ps []*PrimeField

type server struct {
	Params     *PrimeField
	Verifier   *big.Int
	Identifier *big.Int
	varb       *big.Int
	varA       *big.Int
	varB       *big.Int
	varK       *big.Int
	varU       *big.Int
	varM       *big.Int
	Session    *big.Int
	Stage      int
}

func init() {
	ps = PrimeArray
}

func StartServer(vStr string) (*server, error) {
	v, vErr := VerifierFromString(vStr)
	if vErr != nil {
		return nil, vErr
	} else {
		kBi, _ := v.Params.GenK()
		bBI, _ := GenRandom(v.Params.Length / 8)
		serve := &server{
			Params:     v.Params,
			Verifier:   v.Verifier,
			Identifier: v.Identifier,
			varb:       bBI,
			varA:       nil,
			varB:       nil,
			varK:       kBi,
			varU:       nil,
			varM:       nil,
			Session:    nil,
			Stage:      1,
		}
		return serve, nil
	}
}

func (c *server) GenB() (*big.Int, error) {
	if c.Stage == 1 {
		c.varB = big.NewInt(0).Mod(big.NewInt(0).Add(big.NewInt(0).Mul(c.varK, c.Verifier), big.NewInt(0).Exp(c.Params.BG, c.varb, c.Params.BN)), c.Params.BN)
		c.Stage = 2
		return c.varB, nil
	} else {
		return nil, fmt.Errorf("Incorrect stage")
	}
}

func (c *server) SetA(A *big.Int) (bool, error) {
	if c.Stage == 2 {
		c.varA = A
		varU, _ := GenU(c.varA, c.varB)
		c.varU = varU
		c.varM = big.NewInt(0).Exp(big.NewInt(0).Mul(c.varA, big.NewInt(0).Exp(c.Verifier, c.varU, c.Params.BN)), c.varb, c.Params.BN)
		c.Stage = 3
		return true, nil
	} else {
		return false, fmt.Errorf("Incorrect stage")
	}
}

func (c *server) VerifiyH(H string) (bool, error) {
	if c.Stage == 3 {
		_, MStr := GenHash(BI2STR(c.varM))
		if strings.Compare(H, MStr) == 0 {
			sBI, _ := GenHash(MStr)
			c.Session = sBI
			c.Stage = 4
			return true, nil
		} else {
			return false, fmt.Errorf("Incorrect password")
		}
	} else {
		return false, fmt.Errorf("Incorrect stage")
	}
}

func (c *server) VerifiySession(s string, r string) (bool, error) {
	if c.Stage == 4 {
		sStr := BI2STR(c.Session)
		_, SessionStr := GenHash(r + ":" + sStr)
		if strings.Compare(SessionStr, s) == 0 {
			return true, nil
		} else {
			return false, fmt.Errorf("Incorrect Session")
		}
	} else {
		return false, fmt.Errorf("Incorrect stage")
	}
}

func (c *server) ToString() string {
	return c.Params.Converstring() + "-" +
		BI2STR(c.Verifier) + "-" +
		BI2STR(c.Identifier) + "-" +
		BI2STR(c.varb) + "-" +
		BI2STR(c.varA) + "-" +
		BI2STR(c.varB) + "-" +
		BI2STR(c.varK) + "-" +
		BI2STR(c.varU) + "-" +
		BI2STR(c.varM) + "-" +
		BI2STR(c.Session) + "-" +
		strconv.Itoa(c.Stage)
}

func FromString(s string) *server {
	v := strings.Split(s, "-")
	ls, _ := strconv.Atoi(v[10])
	return &server{
		Params:     Str2PF(v[0]),
		Verifier:   STR2BI(v[1]),
		Identifier: STR2BI(v[2]),
		varb:       STR2BI(v[3]),
		varA:       STR2BI(v[4]),
		varB:       STR2BI(v[5]),
		varK:       STR2BI(v[6]),
		varU:       STR2BI(v[7]),
		varM:       STR2BI(v[8]),
		Session:    STR2BI(v[9]),
		Stage:      ls,
	}
}
