package main

import (
	"log"
	"math/big"

	btcutils "github.com/Varunram/essentials/crypto/btc/utils"
)

func ConstructAdaptorSig(x, Px, Py *big.Int, m []byte) (*big.Int,
	*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, *big.Int) {

	t := GetRandomness()
	r := GetRandomness()

	Tx, Ty := Curve.ScalarBaseMult(t.Bytes()) // T = t*G
	Rx, Ry := Curve.ScalarBaseMult(r.Bytes()) // R = r*G

	P := append(Px.Bytes(), Py.Bytes()...)

	RplusTx, RplusTy := Curve.Add(Rx, Ry, Tx, Ty)
	RplusT := append(RplusTx.Bytes(), RplusTy.Bytes()...) // R+T

	HPRTm := btcutils.Sha256(P, RplusT, m)                // H(P||R+T||m)
	HPRTmx := new(big.Int).Mul(BytesToNum(HPRTm), x)      // H(P||R+T||m) * x
	s := new(big.Int).Add(r, new(big.Int).Add(t, HPRTmx)) // s = r + t + H(P||R+T||m) * x

	spr := new(big.Int).Sub(s, t) // s' = s - t (s' is the adaptor signature)
	return spr, r, Rx, Ry, t, Tx, Ty
}

func VerifyAdaptorSig(spr, Rx, Ry, Tx, Ty, Px, Py *big.Int, m []byte) bool {

	sGx, sGy := Curve.ScalarBaseMult(spr.Bytes()) // s*G

	P := append(Px.Bytes(), Py.Bytes()...)

	RplusTx, RplusTy := Curve.Add(Rx, Ry, Tx, Ty)
	RplusT := append(RplusTx.Bytes(), RplusTy.Bytes()...) // R+T

	HPRTm := btcutils.Sha256(P, RplusT, m) // H(P||R+T||m)

	HPRTmPx, HPRTmPy := Curve.ScalarMult(Px, Py, HPRTm) // H(P||R+T||m) * P

	RplusHPRTmPx, RplusHPRTmPy := Curve.Add(Rx, Ry, HPRTmPx, HPRTmPy) // R + H(P||R+T||m) * P
	if sGx.Cmp(RplusHPRTmPx) == 0 && sGy.Cmp(RplusHPRTmPy) == 0 {     // s*G == R + H(P||R+T||m) * P
		return true
	}
	return false
}

func Generate22AdaptorSchnorrChallenge(Jx, Jy, Rax, Ray, Rbx, Rby, Tx, Ty *big.Int, m []byte) *big.Int {
	J := append(Jx.Bytes(), Jy.Bytes()...)

	RARBx, RARBy := Curve.Add(Rax, Ray, Rbx, Rby)

	RARBTx, RARBTy := Curve.Add(RARBx, RARBy, Tx, Ty)
	RARBT := append(RARBTx.Bytes(), RARBTy.Bytes()...)

	HJRARBTm := btcutils.Sha256(J, RARBT, m) // e = H(J || RA+RB+T || m)
	challenge := HJRARBTm
	return BytesToNum(challenge)
}

// https://joinmarket.me/blog/blog/flipping-the-scriptless-script-on-schnorr/
func testadaptorsig() {
	x, Px, Py, err := GetNewKeys()
	if err != nil {
		log.Fatal(err)
	}

	m := []byte("hello world")
	spr, r, Rx, Ry, t, Tx, Ty := ConstructAdaptorSig(x, Px, Py, m)
	log.Println("r=", r, "t=", t)
	if !VerifyAdaptorSig(spr, Rx, Ry, Tx, Ty, Px, Py, m) {
		log.Println("adaptor sigs don't work")
	} else {
		log.Println("adaptor sigs work")
	}
}
