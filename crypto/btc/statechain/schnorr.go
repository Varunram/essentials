package main

import (
	"github.com/pkg/errors"
	"log"
	"math/big"

	btcutils "github.com/Varunram/essentials/crypto/btc/utils"
)

func SchnorrSign(k, Px, Py *big.Int, m []byte, privkey *big.Int) (*big.Int, *big.Int, *big.Int) {

	P := append(Px.Bytes(), Py.Bytes()...)

	Rx, Ry := Curve.ScalarBaseMult(k.Bytes()) // R = k*G
	R := append(Rx.Bytes(), Ry.Bytes()...)

	eByte := btcutils.Sha256(R, P, m)
	e := new(big.Int).SetBytes(eByte)

	sig := new(big.Int).Add(k, new(big.Int).Mul(e, privkey)) // k + hash(R,P,m) * privkey
	return sig, Rx, Ry
}

func SchnorrVerify(sig *big.Int, Rx, Ry *big.Int, Px, Py *big.Int, m []byte) bool {

	P := append(Px.Bytes(), Py.Bytes()...)
	R := append(Rx.Bytes(), Ry.Bytes()...)

	eByte := btcutils.Sha256(R, P, m)
	//e := new(big.Int).SetBytes(eByte)

	// e is a scalar, multiple the scalar with the point P

	ePx, ePy := Curve.ScalarMult(Px, Py, eByte) // H(R,P,m) * P
	cX, cY := Curve.Add(Rx, Ry, ePx, ePy)
	sx, sy := Curve.ScalarBaseMult(sig.Bytes()) // s*G
	if sx.Cmp(cX) == 0 && sy.Cmp(cY) == 0 {
		return true
	}
	return false
}

func Construct22SchnorrPubkey(a, Ax, Ay, b, Bx, By *big.Int) (*big.Int, *big.Int, *big.Int, *big.Int,
	*big.Int, *big.Int, *big.Int, *big.Int) {

	A := append(Ax.Bytes(), Ay.Bytes()...)
	B := append(Bx.Bytes(), By.Bytes()...)

	HAB := btcutils.Sha256(A, B) // H(A||B)

	HHABA := btcutils.Sha256(HAB, A) // H(H(A||B)||A)
	HHABB := btcutils.Sha256(HAB, B) // H(H(A||B)||B)

	Aprx, Apry := Curve.ScalarMult(Ax, Ay, HHABA) // A' = H(H(A||B)||A) * A
	Bprx, Bpry := Curve.ScalarMult(Bx, By, HHABB) // B' = H(H(A||B)||B) * B

	Jx, Jy := Curve.Add(Aprx, Apry, Bprx, Bpry) // J = A'+B'

	apr := new(big.Int).Mul(BytesToNum(HHABA), a) // a' = H(H(A||B)||A) * a
	bpr := new(big.Int).Mul(BytesToNum(HHABB), b) // b' = H(H(A||B)||B) * b

	return Jx, Jy, apr, bpr, Aprx, Apry, Bprx, Bpry
}

func Generate22SchnorrChallenge(Jx, Jy, Rax, Ray, Rbx, Rby *big.Int, m []byte) *big.Int {
	J := append(Jx.Bytes(), Jy.Bytes()...)

	RARBx, RARBy := Curve.Add(Rax, Ray, Rbx, Rby)
	RARB := append(RARBx.Bytes(), RARBy.Bytes()...) // RA + RB
	HJRARBm := btcutils.Sha256(J, RARB, m)          // e = H(J||RA+RB||m)
	challenge := HJRARBm
	return BytesToNum(challenge)
}

func test22Schnorr() {
	a, Ax, Ay, err := GetNewKeys()
	if err != nil {
		log.Fatal(err)
	}
	ra, Rax, Ray, err := GetNewKeys()
	if err != nil {
		log.Fatal(err)
	}

	rb, Rbx, Rby, err := GetNewKeys()
	if err != nil {
		log.Fatal(err)
	}
	b, Bx, By, err := GetNewKeys()
	if err != nil {
		log.Fatal(err)
	}

	m := []byte("hello world")

	Jx, Jy, apr, bpr, _, _, _, _ := Construct22SchnorrPubkey(a, Ax, Ay, b, Bx, By)

	challenge := Generate22SchnorrChallenge(Jx, Jy, Rax, Ray, Rbx, Rby, m)

	sig1 := BlindServerSign(ra, challenge, apr) // ra + challenge*apr
	sig2 := BlindServerSign(rb, challenge, bpr) // rb + challenge*bpr
	sagg := new(big.Int).Add(sig1, sig2)        // ra + rb + challenge(apr + bpr)

	saggx, saggy := Curve.ScalarBaseMult(sagg.Bytes()) // sagg*G

	RARBx, RARBy := Curve.Add(Rax, Ray, Rbx, Rby)           // RA + RB
	eJx, eJy := Curve.ScalarMult(Jx, Jy, challenge.Bytes()) // challenge * J
	RHSx, RHSy := Curve.Add(RARBx, RARBy, eJx, eJy)         // RA+RB + challenge*J

	if saggx.Cmp(RHSx) == 0 && saggy.Cmp(RHSy) == 0 {
		log.Println("22 schnorr works")
	} else {
		log.Fatal("22 schnorr doesn't work")
	}
}

// http://diyhpl.us/wiki/transcripts/building-on-bitcoin/2018/blind-signatures-and-scriptless-scripts/
func testSchnorr() {
	privkey, Px, Py, err := GetNewKeys()
	if err != nil {
		log.Fatal(err)
	}
	//log.Println("PRIVKEY: ", privkey, "PUBKEY: ", pubkey, len(privkey), len(pubkey))
	k := GetRandomness()
	sig, Rx, Ry := SchnorrSign(k, Px, Py, []byte("hello world"), privkey)
	// log.Println("SCHNORR SIG: ", sig)

	if !SchnorrVerify(sig, Rx, Ry, Px, Py, []byte("hello world")) {
		log.Fatal(errors.New("schnorr sigs don't match"))
	} else {
		log.Println("Schnorr signatures work")
	}
}
