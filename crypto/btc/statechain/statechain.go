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

func SchnorrVerify(sig []byte, Rx, Ry *big.Int, Px, Py *big.Int, m []byte) bool {

	P := append(Px.Bytes(), Py.Bytes()...)
	R := append(Rx.Bytes(), Ry.Bytes()...)

	eByte := btcutils.Sha256(append(append(R, P...), m...))
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
	[]byte, []byte, *big.Int, *big.Int, []byte, []byte) {

	alpha := GetRandomness()
	beta := GetRandomness()

	alphaGX, alphaGY := Curve.ScalarBaseMult(alpha)  // alpha*G
	betaPX, betaPY := Curve.ScalarMult(Px, Py, beta) // beta*P

	// need to add Rx, alphax, betapx
	tempX, tempY := Curve.Add(Rx, Ry, alphaGX, alphaGY)   // R + alpha*G
	RprX, RprY := Curve.Add(tempX, tempY, betaPX, betaPY) // R + alpha*G + beta*P

	Rpr := append(RprX.Bytes(), RprY.Bytes()...) // R' = R + alpha*G + beta*P
	P := append(Px.Bytes(), Py.Bytes()...)

	cpr := btcutils.Sha256(append(append(Rpr, P...), m...))  // c' = H(R',P,m)
	c := new(big.Int).Add(BytesToNum(cpr), BytesToNum(beta)) // c = c' + beta

	return alpha, beta, RprX, RprY, cpr, c.Bytes()
}

func BlindServerSign(k *big.Int, cByte []byte, privkey *big.Int) []byte {
	c := BytesToNum(cByte)
	cx := new(big.Int).Mul(c, privkey) // c*x
	sig := new(big.Int).Add(k, cx)     // s = k + c*x
	return sig.Bytes()
}

func BlindClientUnblind(alphaByte []byte, sigByte []byte) []byte {
	alpha := BytesToNum(alphaByte)
	spr := new(big.Int).Add(BytesToNum(sigByte), alpha) // s' = s + alpha
	return spr.Bytes()
}

func MuSig2CreateSign(x1, X1x, X1y, x2, X2x, X2y, r1, R1x, R1y, r2, R2x, R2y *big.Int,
	m []byte) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int) {

	X1 := append(X1x.Bytes(), X1y.Bytes()...)
	X2 := append(X2x.Bytes(), X2y.Bytes()...)

	// L = H(X1,X2)
	L := btcutils.Sha256(append(X1, X2...))

	hash1 := btcutils.Sha256(append(L, X1...)) // H(L,X1)
	hash2 := btcutils.Sha256(append(L, X2...)) // H(L,X2)

	Xx1, Xy1 := Curve.ScalarMult(X1x, X1y, hash1) // H(L,X1)X1
	Xx2, Xy2 := Curve.ScalarMult(X2x, X2y, hash2) // H(L,X2)X2

	Xx, Xy := Curve.Add(Xx1, Xy1, Xx2, Xy2)
	X := append(Xx.Bytes(), Xy.Bytes()...) // X = H(L,X1)X1 + H(L,X2)X2

	Rx, Ry := Curve.Add(R1x, R1y, R2x, R2y)
	R := append(Rx.Bytes(), Ry.Bytes()...) // R = R1 + R2

	HXRm := BytesToNum(btcutils.Sha256(append(X, append(R, m...)...))) // H(X,R,m)
	HLX1 := BytesToNum(btcutils.Sha256(append(L, X1...)))              // H(L,X1)
	HLX2 := BytesToNum(btcutils.Sha256(append(L, X2...)))              // H(L,X2)

	s1 := new(big.Int).Add(r1, new(big.Int).Mul(new(big.Int).Mul(HXRm, HLX1), x1)) // s1 = r1 + H(X,R,m)*H(L,X1)*x1
	s2 := new(big.Int).Add(r2, new(big.Int).Mul(new(big.Int).Mul(HXRm, HLX2), x2)) // s2 = r2+ H(X,R,m)*H(L,X2)*x2
	s := new(big.Int).Add(s1, s2)                                                  // s = s1 + s2

	return Rx, Ry, Xx, Xy, s
}

func MuSig2Verify(Rx, Ry, Xx, Xy, s *big.Int, m []byte) bool {

	sGx, sGy := Curve.ScalarBaseMult(s.Bytes()) // s*G

	X := append(Xx.Bytes(), Xy.Bytes()...)
	R := append(Rx.Bytes(), Ry.Bytes()...)

	HXRm := btcutils.Sha256(append(X, append(R, m...)...)) // H(X,R,m)
	Cx, Cy := Curve.ScalarMult(Xx, Xy, HXRm)               // H(X,R,m)X
	rightX, rightY := Curve.Add(Rx, Ry, Cx, Cy)            // R + H(X,R,m)X

	if sGx.Cmp(rightX) == 0 && sGy.Cmp(rightY) == 0 { // s*G == R + H(X,R,m)X
		return true
	}

	return false
}

func StatechainGenMuSigKey(X1x, X1y, X2x, X2y *big.Int) ([]byte, *big.Int, *big.Int) {

	X1 := append(X1x.Bytes(), X1y.Bytes()...)
	X2 := append(X2x.Bytes(), X2y.Bytes()...)

	// L = H(X1,X2)
	L := btcutils.Sha256(append(X1, X2...))

	hash1 := btcutils.Sha256(append(L, X1...)) // H(L,X1)
	hash2 := btcutils.Sha256(append(L, X2...)) // H(L,X2)

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

func StatechainRequestBlindSig(userSig []byte, blindedMsg []byte, k *big.Int,
	userPubkey []byte, Bx, By *big.Int, nextUserPubkey []byte) ([]byte, error) {
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

	return serverSig, nil
}

func constructEltooTx(address string, privkey string) (string, error) {
	var tx string
	return tx, nil
}

func broadcastTx(tx []byte) error {
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

	if !SchnorrVerify(sig, Rx, Ry, Px, Py, []byte("hello world")) {
		log.Fatal(errors.New("schnorr sigs don't match"))
	} else {
		log.Println("Schnorr signatures work")
	}
}

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
	tx1 := []byte("") // 1 BTC to AX - this stuff must come from the client
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

	broadcastTx(tx1)

	userSig = BlindClientUnblind(alpha, userSig)
	sig = BlindClientUnblind(alpha, sig)

	log.Println("Passing transitory key: ", x, " to C: ", c, "MUSIG PUBKEY: ", AX, "L=", L)
}
