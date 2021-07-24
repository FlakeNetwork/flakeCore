package flakecore

import (
	"testing"
)

func TestFlakeCore(t *testing.T) {
	ver, _ := GenerateVerifier(512, "Test", "Test")
	verStr := ver.ConvertString()

	c, _ := StartClient("Test", "Test", verStr)
	s, _ := StartServer(verStr)

	B, e3 := s.GenB()
	A, e4 := c.GenA()

	x1, e5 := s.SetA(A)
	x2, e6 := c.SetB(B)

	s1, e7 := c.GenH()
	s2, e8 := s.VerifiyH(s1)

	s1, k1, e9 := c.GenSession()
	r1, e10 := s.VerifiySession(s1, k1)

	if s2 != true {
		t.Error("Incorrect Password")
	}

	if r1 != true {
		t.Error("Incorrect Session")
	}

	if (e3 != nil) || (e4 != nil) || (e5 != nil) || (e6 != nil) || (e7 != nil) || (e8 != nil) || (e9 != nil) || (e10 != nil) {
		t.Error("Error in protocol")
	}

	if (x1 != true) || (x2 != true) {
		t.Error("Error in protocol")
	}
}
