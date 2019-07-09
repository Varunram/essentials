package main

import (
	"crypto/rand"
	"github.com/pkg/errors"
	"io"
	"log"
	"math/big"

	btcutils "github.com/Varunram/essentials/crypto/btc/utils"
)

func BlindServerNonce() (*big.Int, *big.Int, *big.Int) {
	k := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, k)
	if err != nil {
		log.Fatal(err)
	}

	Rx, Ry := Curve.ScalarBaseMult(k) // R = k*G
	return BytesToNum(k), Rx, Ry
}

func BlindClientBlind(Rx *big.Int, Ry *big.Int, m []byte, Px, Py *big.Int) (
	*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, *big.Int) {

	alpha := GetRandomness()
	beta := GetRandomness()

	alphaGX, alphaGY := Curve.ScalarBaseMult(alpha.Bytes())  // alpha*G
	betaPX, betaPY := Curve.ScalarMult(Px, Py, beta.Bytes()) // beta*P

	// need to add Rx, alphax, betapx
	tempX, tempY := Curve.Add(Rx, Ry, alphaGX, alphaGY)   // R + alpha*G
	RprX, RprY := Curve.Add(tempX, tempY, betaPX, betaPY) // R + alpha*G + beta*P

	Rpr := append(RprX.Bytes(), RprY.Bytes()...) // R' = R + alpha*G + beta*P
	P := append(Px.Bytes(), Py.Bytes()...)

	cpr := btcutils.Sha256(Rpr, P, m)            // c' = H(R',P,m)
	c := new(big.Int).Add(BytesToNum(cpr), beta) // c = c' + beta

	return alpha, beta, RprX, RprY, BytesToNum(cpr), c
}

func BlindServerSign(k *big.Int, c *big.Int, privkey *big.Int) *big.Int {
	cx := new(big.Int).Mul(c, privkey) // c*x
	sig := new(big.Int).Add(k, cx)     // s = k + c*x
	return sig
}

func BlindClientUnblind(alpha *big.Int, sig *big.Int) *big.Int {
	spr := new(big.Int).Add(sig, alpha) // s' = s + alpha
	return spr
}

// http://diyhpl.us/wiki/transcripts/building-on-bitcoin/2018/blind-signatures-and-scriptless-scripts/
func testBlindSchnorr() {
	privkey, Px, Py, err := GetNewKeys()
	if err != nil {
		log.Fatal(err)
	}
	//log.Println("PRIVKEY: ", privkey, "PUBKEY: ", pubkey, len(privkey), len(pubkey))
	k, Rx, Ry := BlindServerNonce()

	alpha, _, RprX, RprY, _, c := BlindClientBlind(Rx, Ry, []byte("hello world"), Px, Py)
	//log.Println("ALPHA: ", alpha, "BETA: ", beta, "RprX: ", RprX, "RprY: ", RprY, "cpr: ", cpr, "c: ", c)

	blindSig := BlindServerSign(k, c, privkey)
	spr := BlindClientUnblind(alpha, blindSig)

	if !SchnorrVerify(spr, RprX, RprY, Px, Py, []byte("hello world")) {
		log.Fatal(errors.New("blind schnorr sigs don't match"))
	} else {
		log.Println("Blind Schnorr signatures work")
	}

}
