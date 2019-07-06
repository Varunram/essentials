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

var Storage map[string][]byte
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

func SchnorrGetK() []byte {
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

	rightX, rightY := Curve.Add(Rx, Ry, ePx, ePy)

	sx, sy := Curve.ScalarBaseMult(sig) // s*G

	if sx.Cmp(rightX) == 0 && sy.Cmp(rightY) == 0 {
		return true
	}
	return false
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
	Storage = make(map[string][]byte)
}

func GetNewKeys() ([]byte, *big.Int, *big.Int, error) {

	x, err := NewPrivateKey()
	if err != nil {
		return nil, nil, nil, err
	}

	pkx, pky := PubkeyPointsFromPrivkey(x)
	return x.Bytes(), pkx, pky, nil
}

func requestNewPubkey(userPubkey string) ([]byte, error) {
	log.Println("USER: ", userPubkey, "is requesting a new server pubkey")
	var serverPubkey []byte
	var err error

	_, serverPubkeyX, serverPubkeyY, err := GetNewKeys()
	if err != nil {
		return nil, errors.Wrap(err, "could not generate new pubkey, quitting")
	}

	Storage[userPubkey] = SerializeCompressed(serverPubkeyX, serverPubkeyY)
	return serverPubkey, nil
}

func requestBlingSig(userSig string, blindedMsg string, nextUserPubkey string) (string, error) {
	var blindSig string

	var serverPubkey []byte
	Storage[nextUserPubkey] = serverPubkey
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

func main() {
	privkey, pubkeyX, pubkeyY, err := GetNewKeys()
	if err != nil {
		log.Fatal(err)
	}
	//log.Println("PRIVKEY: ", privkey, "PUBKEY: ", pubkey, len(privkey), len(pubkey))
	k := SchnorrGetK()
	sig, Rx, Ry := SchnorrSign(k, pubkeyX, pubkeyY, "hello world", BytesToNum(privkey))
	log.Println("SCHNORR SIG: ", sig)

	if !SchnorrVerify(sig, Rx, Ry, pubkeyX, pubkeyY, "hello world") {
		log.Fatal(errors.New("schnorr sigs don't match"))
	} else {
		log.Println("Schnorr signatures work")
	}
}
