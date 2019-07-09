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

func GetRandomness() *big.Int {
	k := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, k)
	if err != nil {
		log.Fatal(err)
	}

	return BytesToNum(k)
}

func StatechainGenMuSigKey(X1x, X1y, X2x, X2y *big.Int) ([]byte, *big.Int, *big.Int) {

	X1 := append(X1x.Bytes(), X1y.Bytes()...)
	X2 := append(X2x.Bytes(), X2y.Bytes()...)

	// L = H(X1,X2)
	L := btcutils.Sha256(X1, X2)

	hash1 := btcutils.Sha256(L, X1) // H(L,X1)
	hash2 := btcutils.Sha256(L, X2) // H(L,X2)

	Xx1, Xy1 := Curve.ScalarMult(X1x, X1y, hash1) // H(L,X1)X1
	Xx2, Xy2 := Curve.ScalarMult(X2x, X2y, hash2) // H(L,X2)X2

	Xx, Xy := Curve.Add(Xx1, Xy1, Xx2, Xy2) // X = H(L,X1)X1 + H(L,X2)X2

	return L, Xx, Xy
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

func StateServerRequestNewPubkey(userPubkey []byte) (*big.Int, *big.Int, error) {
	privkey, pkX, pkY, err := GetNewKeys()
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not generate new pubkey, quitting")
	}

	// store private key, pkX, pkY for blind signing later when requested
	Storage[string(userPubkey)] = make([]*big.Int, 3)
	Storage[string(userPubkey)][0] = privkey
	Storage[string(userPubkey)][1] = pkX
	Storage[string(userPubkey)][2] = pkY

	return pkX, pkY, nil
}

func StatechainRequestBlindSig(userSig *big.Int, blindedMsg *big.Int, k *big.Int,
	userPubkey []byte, Bx, By *big.Int, nextUserPubkey []byte) (*big.Int, error) {
	//var serverPubkey []byte
	//Storage[nextUserPubkey] = serverPubkey

	// first, lets retrieve the private key associated with the user's pubkey
	var privkey *big.Int

	val, exists := Storage[string(userPubkey)]
	if !exists {
		return nil, errors.New("private key not found in storage")
	}

	privkey = val[0]
	serverSig := BlindServerSign(k, blindedMsg, privkey) // user has signed over the blind message tx2,

	Storage[string(nextUserPubkey)] = make([]*big.Int, 3)

	Storage[string(nextUserPubkey)][0] = privkey
	Storage[string(nextUserPubkey)][1] = Bx
	Storage[string(nextUserPubkey)][2] = By

	Storage[string(userPubkey)][0] = nil
	Storage[string(userPubkey)][1] = new(big.Int)
	Storage[string(userPubkey)][2] = new(big.Int)

	log.Println("USER SIG: ", userSig)
	return serverSig, nil
}

// https://lists.linuxfoundation.org/pipermail/bitcoin-dev/2019-June/017005.html
// https://github.com/RubenSomsen/rubensomsen.github.io/blob/master/img/statechains.pdf
func teststatechain() {
	InitStorage()
	b, Bx, By, err := GetNewKeys()
	if err != nil {
		log.Fatal(err)
	}

	B := SerializeCompressed(Bx, By)
	Ax, Ay, err := StateServerRequestNewPubkey(B) // request A = a*G from the server, A is the server pubkey
	if err != nil {
		log.Fatal(err)
	}

	x, Xx, Xy, err := GetNewKeys() // generate transitory keypair X (x, Xx, Xy)
	if err != nil {
		log.Fatal(err)
	}

	L, AXx, AXy := StatechainGenMuSigKey(Ax, Ay, Xx, Xy)
	// log.Println("L=", L, "AX=", len(AXx.Bytes()), len(AXy.Bytes()))
	// 	tx1 := []byte("") // 1 BTC to AX - this stuff must come from the client
	tx2 := []byte("") // eltoo tx assigning 1 btc back to B - this stuff must come from the client
	m := tx2

	c, Cx, Cy, err := GetNewKeys()
	if err != nil {
		log.Fatal(err)
	}

	C := SerializeCompressed(Cx, Cy)
	AX := SerializeCompressed(AXx, AXy)
	nextUserPubkey := C

	k, Rx, Ry := BlindServerNonce() // generate nonce for signing

	alpha, _, RprX, RprY, _, challenge := BlindClientBlind(Rx, Ry, m, Bx, By) // blind the message and generate challenge
	userSig := BlindServerSign(k, challenge, b)                               // user has signed over the blind message tx2,
	userSig = BlindClientUnblind(alpha, userSig)
	// pass to server the challenge and the userSig so it can add userSig to its sign and return
	// final MuSig tx
	log.Println("USERSIG: ", userSig)
	if !SchnorrVerify(userSig, RprX, RprY, Bx, By, m) {
		log.Fatal("user signature not verified, quitting")
	}

	alpha, _, RprX, RprY, _, challenge = BlindClientBlind(Rx, Ry, m, Ax, Ay) // blind the message and generate challenge
	sig, err := StatechainRequestBlindSig(userSig, challenge, k, B, Bx, By, nextUserPubkey)
	if err != nil {
		log.Fatal(err)
	}

	sig = BlindClientUnblind(alpha, sig)
	log.Println("SERVERSIG:", sig)
	if !SchnorrVerify(sig, RprX, RprY, Ax, Ay, m) {
		log.Fatal("server sig doesn't match, quitting")
	}

	//broadcastTx(tx1)

	userSig = BlindClientUnblind(alpha, userSig)
	sig = BlindClientUnblind(alpha, sig)

	log.Println("Passing transitory key: ", x, " to C: ", c, "MUSIG PUBKEY: ", AX, "L=", L)
}

func main() {
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
	t, Tx, Ty, err := GetNewKeys()
	if err != nil {
		log.Fatal(err)
	}

	m := []byte("hello world")
	Jx, Jy, apr, bpr, Aprx, Apry, Bprx, Bpry := Construct22SchnorrPubkey(a, Ax, Ay, b, Bx, By)
	challenge := Generate22AdaptorSchnorrChallenge(Jx, Jy, Rax, Ray, Rbx, Rby, Tx, Ty, m)

	sig1 := BlindServerSign(rb, challenge, bpr) // rb + challenge*bpr

	sprGx, sprGy := Curve.ScalarBaseMult(sig1.Bytes()) // sig1*G

	eBx, eBy := Curve.ScalarMult(Bprx, Bpry, challenge.Bytes()) // challenge*B
	RHSx, RHSy := Curve.Add(Rbx, Rby, eBx, eBy)                 // R B+ challenge*B

	if sprGx.Cmp(RHSx) == 0 && sprGy.Cmp(RHSy) == 0 { // sig1*G == RB + challenge*B
		log.Println("can verify Bob's adaptor sig")
	} else {
		log.Fatal("can't verify Bob's adaptor sig")
	}

	sig2 := BlindServerSign(ra, challenge, apr) // ra + challenge * apr

	sprGx, sprGy = Curve.ScalarBaseMult(sig2.Bytes()) // sig2*G

	eAx, eAy := Curve.ScalarMult(Aprx, Apry, challenge.Bytes()) // challenge*A
	RHSx, RHSy = Curve.Add(Rax, Ray, eAx, eAy)                  // RA + challenge*A

	if sprGx.Cmp(RHSx) == 0 && sprGy.Cmp(RHSy) == 0 { // sig2*G == RA+challenge*A
		log.Println("can verify Alice's adaptor sig")
	} else {
		log.Fatal("can't verify Alice's adaptor sig")
	}

	// bob has alice's signature
	rbt := new(big.Int).Add(rb, t)           // rb + t
	ebpr := new(big.Int).Mul(challenge, bpr) // challenge * bpr

	sagg := new(big.Int).Add(sig2, new(big.Int).Add(rbt, ebpr)) // sig2 + rb + t + challenge * bpr

	// alice wants to check the sig
	check := new(big.Int).Sub(new(big.Int).Sub(sagg, sig2), sig1) // sagg - sig2 - sig1 = t

	if check.Cmp(t) == 0 { // check == t?
		log.Println("22 Adaptor Schnorr works")
	} else {
		log.Fatal("22 Adaptor Schnorr doesn't work")
	}
}
