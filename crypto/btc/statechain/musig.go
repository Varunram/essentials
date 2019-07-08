package main

import (
	"log"
	"math/big"

	btcutils "github.com/Varunram/essentials/crypto/btc/utils"
)

func MuSig2CreateSign(x1, X1x, X1y, x2, X2x, X2y, r1, R1x, R1y, r2, R2x, R2y *big.Int,
	m []byte) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int) {

	X1 := append(X1x.Bytes(), X1y.Bytes()...)
	X2 := append(X2x.Bytes(), X2y.Bytes()...)

	// L = H(X1,X2)
	L := btcutils.Sha256(X1, X2)

	hash1 := btcutils.Sha256(L, X1) // H(L,X1)
	hash2 := btcutils.Sha256(L, X2) // H(L,X2)

	Xx1, Xy1 := Curve.ScalarMult(X1x, X1y, hash1) // H(L,X1)X1
	Xx2, Xy2 := Curve.ScalarMult(X2x, X2y, hash2) // H(L,X2)X2

	Xx, Xy := Curve.Add(Xx1, Xy1, Xx2, Xy2)
	X := append(Xx.Bytes(), Xy.Bytes()...) // X = H(L,X1)X1 + H(L,X2)X2

	Rx, Ry := Curve.Add(R1x, R1y, R2x, R2y)
	R := append(Rx.Bytes(), Ry.Bytes()...) // R = R1 + R2

	HXRm := BytesToNum(btcutils.Sha256(X, R, m)) // H(X,R,m)
	HLX1 := BytesToNum(btcutils.Sha256(L, X1))   // H(L,X1)
	HLX2 := BytesToNum(btcutils.Sha256(L, X2))   // H(L,X2)

	s1 := new(big.Int).Add(r1, new(big.Int).Mul(new(big.Int).Mul(HXRm, HLX1), x1)) // s1 = r1 + H(X,R,m)*H(L,X1)*x1
	s2 := new(big.Int).Add(r2, new(big.Int).Mul(new(big.Int).Mul(HXRm, HLX2), x2)) // s2 = r2+ H(X,R,m)*H(L,X2)*x2
	s := new(big.Int).Add(s1, s2)                                                  // s = s1 + s2

	return Rx, Ry, Xx, Xy, s
}

func MuSig2Verify(Rx, Ry, Xx, Xy, s *big.Int, m []byte) bool {

	sGx, sGy := Curve.ScalarBaseMult(s.Bytes()) // s*G

	X := append(Xx.Bytes(), Xy.Bytes()...)
	R := append(Rx.Bytes(), Ry.Bytes()...)

	HXRm := btcutils.Sha256(X, R, m)            // H(X,R,m)
	Cx, Cy := Curve.ScalarMult(Xx, Xy, HXRm)    // H(X,R,m)X
	rightX, rightY := Curve.Add(Rx, Ry, Cx, Cy) // R + H(X,R,m)X

	if sGx.Cmp(rightX) == 0 && sGy.Cmp(rightY) == 0 { // s*G == R + H(X,R,m)X
		return true
	}

	return false
}

// https://blockstream.com/2018/01/23/en-musig-key-aggregation-schnorr-signatures/
func testmusig() {
	p1, P1x, P1y, err := GetNewKeys()
	if err != nil {
		log.Fatal(err)
	}

	p2, P2x, P2y, err := GetNewKeys()
	if err != nil {
		log.Fatal(err)
	}

	r1, R1x, R1y := BlindServerNonce() // craete random ri
	r2, R2x, R2y := BlindServerNonce() // craete random ri

	message := []byte("hello world")

	Rx, Ry, Xx, Xy, s := MuSig2CreateSign(p1, P1x, P1y, p2, P2x, P2y, r1, R1x, R1y, r2, R2x, R2y, message)
	if MuSig2Verify(Rx, Ry, Xx, Xy, s, message) {
		log.Println("musig verify works")
	} else {
		log.Println("musig verify doesn't work")
	}
}
