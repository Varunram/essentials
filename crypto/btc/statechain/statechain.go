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

func SchnorrSign(kByte []byte, Px, Py *big.Int, m []byte, privkey *big.Int) (*big.Int, *big.Int, *big.Int) {

	P := append(Px.Bytes(), Py.Bytes()...)

	Rx, Ry := Curve.ScalarBaseMult(kByte) // R = k*G
	R := append(Rx.Bytes(), Ry.Bytes()...)

	eByte := btcutils.Sha256(R, P, m)
	e := new(big.Int).SetBytes(eByte)

	k := new(big.Int).SetBytes(kByte) // hash(R,P,m)

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

	cpr := btcutils.Sha256(Rpr, P, m)                        // c' = H(R',P,m)
	c := new(big.Int).Add(BytesToNum(cpr), BytesToNum(beta)) // c = c' + beta

	return alpha, beta, RprX, RprY, cpr, c.Bytes()
}

func BlindServerSign(k *big.Int, cByte []byte, privkey *big.Int) *big.Int {
	c := BytesToNum(cByte)
	cx := new(big.Int).Mul(c, privkey) // c*x
	sig := new(big.Int).Add(k, cx)     // s = k + c*x
	return sig
}

func BlindClientUnblind(alphaByte []byte, sig *big.Int) *big.Int {
	alpha := BytesToNum(alphaByte)
	spr := new(big.Int).Add(sig, alpha) // s' = s + alpha
	return spr
}

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

func StatechainRequestBlindSig(userSig *big.Int, blindedMsg []byte, k *big.Int,
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

func ConstructAdaptorSig(x, Px, Py *big.Int, m []byte) (*big.Int,
	*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, *big.Int) {

	t := GetRandomness()
	r := GetRandomness()

	Tx, Ty := Curve.ScalarBaseMult(t) // T = t*G
	Rx, Ry := Curve.ScalarBaseMult(r) // R = r*G

	P := append(Px.Bytes(), Py.Bytes()...)

	RplusTx, RplusTy := Curve.Add(Rx, Ry, Tx, Ty)
	RplusT := append(RplusTx.Bytes(), RplusTy.Bytes()...) // R+T

	HPRTm := btcutils.Sha256(P, RplusT, m)                                        // H(P||R+T||m)
	HPRTmx := new(big.Int).Mul(BytesToNum(HPRTm), x)                              // H(P||R+T||m) * x
	s := new(big.Int).Add(BytesToNum(r), new(big.Int).Add(BytesToNum(t), HPRTmx)) // s = r + t + H(P||R+T||m) * x

	spr := new(big.Int).Sub(s, BytesToNum(t)) // s' = s - t (s' is the adaptor signature)
	return spr, BytesToNum(r), Rx, Ry, BytesToNum(t), Tx, Ty
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

func Generate22SchnorrChallenge(Jx, Jy, Rax, Ray, Rbx, Rby *big.Int, m []byte) []byte {
	J := append(Jx.Bytes(), Jy.Bytes()...)

	RARBx, RARBy := Curve.Add(Rax, Ray, Rbx, Rby)
	RARB := append(RARBx.Bytes(), RARBy.Bytes()...) // RA + RB
	HJRARBm := btcutils.Sha256(J, RARB, m)          // e = H(J||RA+RB||m)
	challenge := HJRARBm
	return challenge
}

func Generate22AdaptorSchnorrChallenge(Jx, Jy, Rax, Ray, Rbx, Rby, Tx, Ty *big.Int, m []byte) []byte {
	J := append(Jx.Bytes(), Jy.Bytes()...)

	RARBx, RARBy := Curve.Add(Rax, Ray, Rbx, Rby)

	RARBTx, RARBTy := Curve.Add(RARBx, RARBy, Tx, Ty)
	RARBT := append(RARBTx.Bytes(), RARBTy.Bytes()...)

	HJRARBTm := btcutils.Sha256(J, RARBT, m) // e = H(J || RA+RB+T || m)
	challenge := HJRARBTm
	return challenge
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

	RARBx, RARBy := Curve.Add(Rax, Ray, Rbx, Rby)   // RA + RB
	eJx, eJy := Curve.ScalarMult(Jx, Jy, challenge) // challenge * J
	RHSx, RHSy := Curve.Add(RARBx, RARBy, eJx, eJy) // RA+RB + challenge*J

	if saggx.Cmp(RHSx) == 0 && saggy.Cmp(RHSy) == 0 {
		log.Println("22 schnorr works")
	} else {
		log.Fatal("22 schnorr doesn't work")
	}
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

	eBx, eBy := Curve.ScalarMult(Bprx, Bpry, challenge) // challenge*B
	RHSx, RHSy := Curve.Add(Rbx, Rby, eBx, eBy)         // R B+ challenge*B

	if sprGx.Cmp(RHSx) == 0 && sprGy.Cmp(RHSy) == 0 { // sig1*G == RB + challenge*B
		log.Println("can verify Bob's adaptor sig")
	} else {
		log.Fatal("can't verify Bob's adaptor sig")
	}

	sig2 := BlindServerSign(ra, challenge, apr) // ra + challenge * apr

	sprGx, sprGy = Curve.ScalarBaseMult(sig2.Bytes()) // sig2*G

	eAx, eAy := Curve.ScalarMult(Aprx, Apry, challenge) // challenge*A
	RHSx, RHSy = Curve.Add(Rax, Ray, eAx, eAy)          // RA + challenge*A

	if sprGx.Cmp(RHSx) == 0 && sprGy.Cmp(RHSy) == 0 { // sig2*G == RA+challenge*A
		log.Println("can verify Alice's adaptor sig")
	} else {
		log.Fatal("can't verify Alice's adaptor sig")
	}

	// bob has alice's signature
	rbt := new(big.Int).Add(rb, t)                       // rb + t
	ebpr := new(big.Int).Mul(BytesToNum(challenge), bpr) // challenge * bpr

	sagg := new(big.Int).Add(sig2, new(big.Int).Add(rbt, ebpr)) // sig2 + rb + t + challenge * bpr

	// alice wants to check the sig
	check := new(big.Int).Sub(new(big.Int).Sub(sagg, sig2), sig1) // sagg - sig2 - sig1 = t

	if check.Cmp(t) == 0 { // check == t?
		log.Println("22 Adaptor Schnorr works")
	} else {
		log.Fatal("22 Adaptor Schnorr doesn't work")
	}
}
