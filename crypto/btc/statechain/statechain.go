package main

import (
	"crypto/rand"
	"github.com/pkg/errors"
	"io"
	"log"
	"math/big"

	btcutils "github.com/Varunram/essentials/crypto/btc/utils"
	"github.com/btcsuite/btcd/btcec"
)

var Storage map[string][]*big.Int
var Curve *btcec.KoblitzCurve = btcec.S256() // take only the curve, can't use other stuff

func NewPrivateKey() (*big.Int, error) {
	b := make([]byte, Curve.Params().BitSize/8+8)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		log.Fatal(err)
	}

	var one = new(big.Int).SetInt64(1)
	x := new(big.Int).SetBytes(b)
	n := new(big.Int).Sub(Curve.Params().N, one)
	x.Mod(x, n)
	x.Add(x, one)

	return x, nil
}

func PubkeyPointsFromPrivkey(privkey *big.Int) (*big.Int, *big.Int) {
	x, y := Curve.ScalarBaseMult(privkey.Bytes())
	return x, y
}

func BytesToNum(byteString []byte) *big.Int {
	return new(big.Int).SetBytes(byteString)
}

func GetRandomness() []byte {
	k := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, k)
	if err != nil {
		log.Fatal(err)
	}

	return k
}

func SchnorrSign(kByte []byte, Px, Py *big.Int, m string, privkey *big.Int) ([]byte, *big.Int, *big.Int) {

	P := append(Px.Bytes(), Py.Bytes()...)

	Rx, Ry := Curve.ScalarBaseMult(kByte) // R = k*G
	R := append(Rx.Bytes(), Ry.Bytes()...)

	eByte := btcutils.Sha256(append(append(R, P...), []byte(m)...))
	e := new(big.Int).SetBytes(eByte)

	k := new(big.Int).SetBytes(kByte) // hash(R,P,m)

	sig := new(big.Int).Add(k, new(big.Int).Mul(e, privkey)) // k + hash(R,P,m) * privkey
	return sig.Bytes(), Rx, Ry
}

func SchnorrVerify(sig []byte, Rx, Ry *big.Int, Px, Py *big.Int, m string) bool {

	P := append(Px.Bytes(), Py.Bytes()...)
	R := append(Rx.Bytes(), Ry.Bytes()...)

	eByte := btcutils.Sha256(append(append(R, P...), []byte(m)...))
	//e := new(big.Int).SetBytes(eByte)

	// e is a scalar, multiple the scalar with the point P

	ePx, ePy := Curve.ScalarMult(Px, Py, eByte) // H(R,P,m) * P
	cX, cY := Curve.Add(Rx, Ry, ePx, ePy)
	sx, sy := Curve.ScalarBaseMult(sig) // s*G
	if sx.Cmp(cX) == 0 && sy.Cmp(cY) == 0 {
		return true
	}
	return false
}

func BlindServerNonce() ([]byte, *big.Int, *big.Int) {
	k := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, k)
	if err != nil {
		log.Fatal(err)
	}

	Rx, Ry := Curve.ScalarBaseMult(k) // R = k*G
	return k, Rx, Ry
}

func BlindClientBlind(Rx *big.Int, Ry *big.Int, m string, Px, Py *big.Int, privkey *big.Int) (
	[]byte, []byte, *big.Int, *big.Int, []byte, []byte) {
	alpha := GetRandomness()
	beta := GetRandomness()

	alphaGX, alphaGY := Curve.ScalarBaseMult(alpha)
	betaPX, betaPY := Curve.ScalarMult(Px, Py, beta)

	// need to add Rx, alphax, betapx
	tempX, tempY := Curve.Add(Rx, Ry, alphaGX, alphaGY) // R + alpha*G
	RprX, RprY := Curve.Add(tempX, tempY, betaPX, betaPY) // R + alpha*G + beta*P

	Rpr := append(RprX.Bytes(), RprY.Bytes()...)
	P := append(Px.Bytes(), Py.Bytes()...)

	cpr := btcutils.Sha256(append(append(Rpr, P...), []byte(m)...))
	c := new(big.Int).Add(BytesToNum(cpr), BytesToNum(beta))

	return alpha, beta, RprX, RprY, cpr, c.Bytes()
}

func BlindServerSign(kByte []byte, cByte []byte, privkey *big.Int) []byte {
	k := BytesToNum(kByte)
	c := BytesToNum(cByte)
	cx := new(big.Int).Mul(c, privkey)

	sig := new(big.Int).Add(k, cx)
	return sig.Bytes()
}

func BlindClientUnblind(alphaByte []byte, sigByte []byte) []byte {
	alpha := BytesToNum(alphaByte)
	s := BytesToNum(sigByte)
	spr := new(big.Int).Add(s, alpha)
	return spr.Bytes()
}

// SerializeCompressed serializes a public key in a 33-byte compressed format.
func SerializeCompressed(pkx *big.Int, pky *big.Int) []byte {
	b := make([]byte, 0, 33)
	format := byte(0x02) // magic number for ybyte + xcoord
	if pky.Bit(0) == 1 {
		format |= 0x1
	}
	b = append(b, format)
	return append(b, pkx.Bytes()...)
}

func InitStorage() {
	Storage = make(map[string][]*big.Int)
}

func GetNewKeys() (*big.Int, *big.Int, *big.Int, error) {

	x, err := NewPrivateKey()
	if err != nil {
		return nil, nil, nil, err
	}

	pkx, pky := PubkeyPointsFromPrivkey(x)
	return x, pkx, pky, nil
}

func StateServerRequestNewPubkey(userPubkey string) ([]byte, error) {
	log.Println("USER: ", userPubkey, "is requesting a new server pubkey")
	var serverPubkey []byte
	var err error

	privkey, pkX, pkY, err := GetNewKeys()
	if err != nil {
		return nil, errors.Wrap(err, "could not generate new pubkey, quitting")
	}

	// store private key, pkX, pkY for blind signing later when requested
	Storage[userPubkey][0] = privkey
	Storage[userPubkey][1] = pkX
	Storage[userPubkey][2] = pkY

	return serverPubkey, nil
}

func requestBlindSig(userSig string, blindedMsg string, nextUserPubkey string) (string, error) {
	var blindSig string

	//var serverPubkey []byte
	//Storage[nextUserPubkey] = serverPubkey
	return blindSig, nil
}

func genTransitoryKey() ([]byte, []byte, error) {
	privkey, err := btcec.NewPrivateKey(Curve)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not get new privkey, quitting")
	}

	pubkey := privkey.PubKey().SerializeCompressed()
	return nil, pubkey, nil
}

func getMuSigKey(a, A, x, X string) (string, error) {
	return a + x, nil
}

func constructMusigTx(amount string, address string, privkey string) (string, error) {
	var tx string
	return tx, nil
}

func constructEltooTx(address string, privkey string) (string, error) {
	var tx string
	return tx, nil
}

func broadcastTx(tx string) error {
	return nil
}

func testSchnorr() {
	privkey, Px, Py, err := GetNewKeys()
	if err != nil {
		log.Fatal(err)
	}
	//log.Println("PRIVKEY: ", privkey, "PUBKEY: ", pubkey, len(privkey), len(pubkey))
	k := GetRandomness()
	sig, Rx, Ry := SchnorrSign(k, Px, Py, "hello world", privkey)
	// log.Println("SCHNORR SIG: ", sig)

	if !SchnorrVerify(sig, Rx, Ry, Px, Py, "hello world") {
		log.Fatal(errors.New("schnorr sigs don't match"))
	} else {
		log.Println("Schnorr signatures work")
	}
}

func main() {
	privkey, Px, Py, err := GetNewKeys()
	if err != nil {
		log.Fatal(err)
	}
	//log.Println("PRIVKEY: ", privkey, "PUBKEY: ", pubkey, len(privkey), len(pubkey))
	k, Rx, Ry := BlindServerNonce()

	alpha, _, RprX, RprY, _, c := BlindClientBlind(Rx, Ry, "hello world", Px, Py, privkey)
	//log.Println("ALPHA: ", alpha, "BETA: ", beta, "RprX: ", RprX, "RprY: ", RprY, "cpr: ", cpr, "c: ", c)

	blindSig := BlindServerSign(k, c, privkey)
	spr := BlindClientUnblind(alpha, blindSig)

	if !SchnorrVerify(spr, RprX, RprY, Px, Py, "hello world") {
		log.Fatal(errors.New("blind schnorr sigs don't match"))
	} else {
		log.Println("Blind Schnorr signatures work")
	}

}
