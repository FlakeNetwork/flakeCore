package flakecore

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

var s128 = "128:16:2:18D9180FAA15126BCB5AD2BC6722EAAFB"
var s256 = "256:32:2:1C5A340BABF6A7C902F7B415DF7DBE42236905598C1884F4DE3F2109F0E48ECC3"
var s512 = "512:64:7:1B146F489A0B9EE3686BD715E9B3DD73937E1E2D7B8BA2547860C884EE413BF2BC4529EFDBFF47B566F5B263032FDC34197F1B0137444319821F2E156803CD84F"

var pfStr = []*string{
	&s128,
	&s256,
	&s512,
}

type Verifier struct {
	Params     *PrimeField
	Identifier *big.Int
	Verifier   *big.Int
	Salt       *big.Int
}

var PrimeArray []*PrimeField

func (v *Verifier) ConvertString() string {
	idStr := BI2STR(v.Identifier)
	salt := BI2STR(v.Salt)
	verifier := BI2STR(v.Verifier)
	return strconv.Itoa(v.Params.Length) + ":" + idStr + ":" + verifier + ":" + salt
}

func BI2STR(n *big.Int) string {
	return fmt.Sprintf("%X", n) // or %X or upper case
}

func STR2BI(s string) *big.Int {
	x := new(big.Int)
	x.SetString(s[0:len(s)], 16)
	return x
}

func GenHash(s string) (*big.Int, string) {
	sbyte := []byte(s)
	sHash := sha256.Sum256(sbyte)
	hexStr := hex.EncodeToString(sHash[:])
	bint := STR2BI(hexStr)
	hexStr = strings.ToUpper(hexStr)
	return bint, hexStr
}

func GenU(A *big.Int, B *big.Int) (*big.Int, string) {
	strA := BI2STR(A)
	strB := BI2STR(B)
	sbyte := []byte(strA + ":" + strB)
	sHash := sha256.Sum256(sbyte)
	hexStr := hex.EncodeToString(sHash[:])
	bint := STR2BI(hexStr)
	hexStr = strings.ToUpper(hexStr)
	return bint, hexStr
}

func GenX(s string, i string, p string) (*big.Int, string) {
	sbyte := []byte(s + ":" + i + ":" + p)
	sHash := sha256.Sum256(sbyte)
	hexStr := hex.EncodeToString(sHash[:])
	bint := STR2BI(hexStr)
	hexStr = strings.ToUpper(hexStr)
	return bint, hexStr
}

func GenRandom(s int) (*big.Int, string) {
	sHash := make([]byte, s)
	rand.Read(sHash)
	hexStr := hex.EncodeToString(sHash[:])
	bint := STR2BI(hexStr)
	hexStr = strings.ToUpper(hexStr)
	return bint, hexStr
}

func VerifierFromString(s string) (*Verifier, error) {
	var state bool
	splitS := strings.Split(s, ":")
	var param *PrimeField
	l, _ := strconv.Atoi(splitS[0])
	for _, i := range PrimeArray {
		if i.Length == l {
			param = i
			state = true
		}
	}
	if !state {
		return nil, fmt.Errorf("Invalid Length")
	}
	ver := &Verifier{
		Params:     param,
		Identifier: STR2BI(splitS[1]),
		Verifier:   STR2BI(splitS[2]),
		Salt:       STR2BI(splitS[3]),
	}
	return ver, nil
}
