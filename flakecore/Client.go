package flakecore

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

var pc []*PrimeField

type client struct {
	Params     *PrimeField
	Salt       *big.Int
	Identifier *big.Int
	varx       *big.Int
	vara       *big.Int
	varA       *big.Int
	varB       *big.Int
	varK       *big.Int
	varU       *big.Int
	varM       *big.Int
	Session    *big.Int
	Stage      int
}

func GenerateVerifier(l int, u string, p string) (*Verifier, error) {
	var state bool
	var verify *Verifier

	for _, i := range pc {
		if i.Length == l {
			state = true
			stBI, stStr := GenRandom(32)
			idBI, _ := GenHash(u)
			xBI, _ := GenX(stStr, u, p)
			v := big.NewInt(0).Exp(i.BG, xBI, i.BN)
			verify = &Verifier{
				Params:     i,
				Identifier: idBI,
				Verifier:   v,
				Salt:       stBI,
			}
		}
	}
	if !state {
		return nil, fmt.Errorf("Invlid Length ")
	}
	return verify, nil
}

func StartClient(u string, p string, vString string) (*client, error) {
	var clientObj *client
	var state bool
	vSplit := strings.Split(vString, ":")
	l, _ := strconv.Atoi(vSplit[0])
	_, iHashStr := GenHash(u)
	xBI, _ := GenX(vSplit[3], u, p)
	for _, i := range pc {
		if (i.Length == l) && (iHashStr == vSplit[1]) {
			state = true
			kBI, _ := i.GenK()
			varaBi, _ := GenRandom(i.Length / 8)
			clientObj = &client{
				Params:     i,
				Salt:       STR2BI(vSplit[3]),
				Identifier: STR2BI(vSplit[1]),
				varx:       xBI,
				vara:       varaBi,
				varA:       big.NewInt(0),
				varB:       big.NewInt(0),
				varK:       kBI,
				varU:       big.NewInt(0),
				varM:       big.NewInt(0),
				Session:    big.NewInt(0),
				Stage:      1,
			}
		}
	}
	if state == false {
		return nil, fmt.Errorf("INvalid Length or Invalid String ")
	}
	return clientObj, nil
}

func (c *client) GenA() (*big.Int, error) {
	if c.Stage == 1 {
		c.varA = big.NewInt(0).Exp(c.Params.BG, c.vara, c.Params.BN)
		c.Stage = 2
		return c.varA, nil
	} else {
		return nil, fmt.Errorf("Invalid Stage")
	}
}

func (c *client) SetB(B *big.Int) (bool, error) {
	if c.Stage == 2 {
		c.varB = B
		varU, _ := GenU(c.varA, c.varB)
		c.varU = varU
		c.varM = big.NewInt(0).Exp(big.NewInt(0).Sub(c.varB, big.NewInt(0).Mul(c.varK, big.NewInt(0).Exp(c.Params.BG, c.varx, c.Params.BN))), big.NewInt(0).Add(c.vara, big.NewInt(0).Mul(c.varU, c.varx)), c.Params.BN)
		c.Stage = 3
		return true, nil
	} else {
		return false, fmt.Errorf("Invalid Stage")
	}
}

func (c *client) FetchM() *big.Int {
	return c.varM
}

func (c *client) GenH() (string, error) {
	if c.Stage == 3 {
		_, mStr := GenHash(BI2STR(c.varM))
		sBI, _ := GenHash(mStr)
		c.Session = sBI
		c.Stage = 4
		return mStr, nil
	} else {
		return "", fmt.Errorf("Invalid Stage")
	}
}

func (c *client) GenSession() (string, string, error) {
	if c.Stage == 4 {
		_, rStr := GenRandom(c.Params.SN)
		sStr := BI2STR(c.Session)
		_, SessionStr := GenHash(rStr + ":" + sStr)
		return SessionStr, rStr, nil
	} else {
		return "", "", fmt.Errorf("Invalid Stage")
	}
}

func (c *client) ToString() string {
	return c.Params.Converstring() + "-" +
		BI2STR(c.Salt) + "-" +
		BI2STR(c.Identifier) + "-" +
		BI2STR(c.varx) + "-" +
		BI2STR(c.vara) + "-" +
		BI2STR(c.varA) + "-" +
		BI2STR(c.varB) + "-" +
		BI2STR(c.varK) + "-" +
		BI2STR(c.varU) + "-" +
		BI2STR(c.varM) + "-" +
		BI2STR(c.Session) + "-" +
		strconv.Itoa(c.Stage)
}

func FromStringClient(s string) *client {
	v := strings.Split(s, "-")
	ls, _ := strconv.Atoi(v[11])
	return &client{
		Params:     Str2PF(v[0]),
		Salt:       STR2BI(v[1]),
		Identifier: STR2BI(v[2]),
		varx:       STR2BI(v[3]),
		vara:       STR2BI(v[4]),
		varA:       STR2BI(v[5]),
		varB:       STR2BI(v[6]),
		varK:       STR2BI(v[7]),
		varU:       STR2BI(v[8]),
		varM:       STR2BI(v[9]),
		Session:    STR2BI(v[10]),
		Stage:      ls,
	}
}
