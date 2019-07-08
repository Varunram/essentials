/*
	Copyright 2013-present wemeetagain https://github.com/wemeetagain/go-hdwallet
	Copyright 2019-present Varunram Ganesh
*/
package utils

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"golang.org/x/crypto/ripemd160"
)

var Curve *btcec.KoblitzCurve = btcec.S256()

func Hash160(data []byte) []byte {
	sha := sha256.New()
	ripe := ripemd160.New()
	sha.Write(data)
	ripe.Write(sha.Sum(nil))
	return ripe.Sum(nil)
}

func Sha256(inputs ...[]byte) []byte {
	shaNew := sha256.New()
	for _, input := range inputs {
		shaNew.Write(input)
	}
	return shaNew.Sum(nil)
}

func DoubleSha256(data []byte) []byte {
	return Sha256(Sha256(data))
}

func PrivToPub(key []byte) []byte {
	return Compress(Curve.ScalarBaseMult(key))
}

func Compress(x, y *big.Int) []byte {
	two := big.NewInt(2)
	rem := two.Mod(y, two).Uint64()
	rem += 2
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(rem))
	rest := x.Bytes()
	pad := 32 - len(rest)
	if pad != 0 {
		zeroes := make([]byte, pad)
		rest = append(zeroes, rest...)
	}
	return append(b[1:], rest...)
}

//2.3.4 of SEC1 - http://www.secg.org/index.php?action=secg,docs_secg
func Expand(key []byte) (*big.Int, *big.Int) {
	params := Curve.Params()
	exp := big.NewInt(1)
	exp.Add(params.P, exp)
	exp.Div(exp, big.NewInt(4))
	x := big.NewInt(0).SetBytes(key[1:33])
	y := big.NewInt(0).SetBytes(key[:1])
	beta := big.NewInt(0)
	beta.Exp(x, big.NewInt(3), nil)
	beta.Add(beta, big.NewInt(7))
	beta.Exp(beta, exp, params.P)
	if y.Add(beta, y).Mod(y, big.NewInt(2)).Int64() == 0 {
		y = beta
	} else {
		y = beta.Sub(params.P, beta)
	}
	return x, y
}

func AddPrivKeys(k1, k2 []byte) []byte {
	i1 := big.NewInt(0).SetBytes(k1)
	i2 := big.NewInt(0).SetBytes(k2)
	i1.Add(i1, i2)
	i1.Mod(i1, Curve.Params().N)
	k := i1.Bytes()
	zero, _ := hex.DecodeString("00")
	return append(zero, k...)
}

func AddPubKeys(k1, k2 []byte) []byte {
	x1, y1 := Expand(k1)
	x2, y2 := Expand(k2)
	return Compress(Curve.Add(x1, y1, x2, y2))
}
