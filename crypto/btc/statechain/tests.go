package main

import (
	"github.com/pkg/errors"
	"log"
)

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

// https://blockstream.com/2018/01/23/en-musig-key-aggregation-schnorr-signatures/
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
