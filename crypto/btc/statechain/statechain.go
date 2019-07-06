package main

import (
	"crypto/rand"
	"github.com/pkg/errors"
	"io"
	"log"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	btcutils "github.com/Varunram/essentials/crypto/btc/utils"
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

func SchnorrGetR() ([]byte, *big.Int, *big.Int) {
	k := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, k)
	if err != nil {
		log.Fatal(err)
	}

	x, y := Curve.ScalarBaseMult(k) // R = k*G
	return k, x, y
}

func SchnorrSign(k []byte, Rx, Ry *big.Int, P []byte, m string, privkey *big.Int) {
	R := append(Rx.Bytes(), Ry.Bytes()...)

	e := btcutils.Sha256(append(append(R, P...), []byte(m)...))

	eNum := new(big.Int).SetBytes(e)
	kNum := new(big.Int).SetBytes(k)

	sig := new(big.Int).Add(kNum, new(big.Int).Mul(eNum, privkey)) // k + hash(R,P,m) * privkey
	log.Println("SIG: ", sig.Bytes(), len(sig.Bytes()))
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

func GetNewKeys() ([]byte, []byte, error) {

	x, err := NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}

	pkx, pky := PubkeyPointsFromPrivkey(x)
	return x.Bytes(), SerializeCompressed(pkx, pky), nil
}

func requestNewPubkey(userPubkey string) ([]byte, error) {
	log.Println("USER: ", userPubkey, "is requesting a new server pubkey")
	var serverPubkey []byte
	var err error

	_, serverPubkey, err = GetNewKeys()
	if err != nil {
		return nil, errors.Wrap(err, "could not generate new pubkey, quitting")
	}

	Storage[userPubkey] = serverPubkey
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
	privkey, pubkey, err := GetNewKeys()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("PRIVKEY: ", privkey, "PUBKEY: ", pubkey, len(privkey), len(pubkey))

	k, Rx, Ry := SchnorrGetR()

	SchnorrSign(k, Rx, Ry, pubkey, "hello world", BytesToNum(privkey))

}
