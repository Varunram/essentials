package hdwallet

import (
	"encoding/hex"
	"github.com/pkg/errors"
	"log"
	"testing"
)

// see https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki#test-vectors for test vectors
func testHelper(privkey *HDWallet, matchPrivKey string, pubkey *HDWallet, matchPubkey string) error {
	if privkey.String() != matchPrivKey {
		log.Println(privkey.String())
		return errors.New("private key doesn't match")
	}

	if pubkey.String() != matchPubkey {
		log.Println(matchPubkey)
		return errors.New("public key doesn't match")
	}
	return nil
}

func TestVector1Bip32(t *testing.T) {
	seed, err := hex.DecodeString("000102030405060708090a0b0c0d0e0f")
	if err != nil {
		t.Fatal(err)
	}

	masterprv := MasterKey(seed)
	if masterprv.String() != "xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi" {
		t.Fatal(errors.New("master priv key doesn't match, quitting"))
	}

	masterpub := masterprv.Pub()
	if masterpub.String() != "xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EGMcet8" {
		t.Fatal(errors.New("master pub key doesn't match, quitting"))
	}

	hKey := uint32(0x80000000)
	chainPrivm0H, err := masterprv.Child(hKey)
	if err != nil {
		t.Fatal(err)
	}
	chainPubm0H := chainPrivm0H.Pub()

	err = testHelper(chainPrivm0H, "xprv9uHRZZhk6KAJC1avXpDAp4MDc3sQKNxDiPvvkX8Br5ngLNv1TxvUxt4cV1rGL5hj6KCesnDYUhd7oWgT11eZG7XnxHrnYeSvkzY7d2bhkJ7", chainPubm0H, "xpub68Gmy5EdvgibQVfPdqkBBCHxA5htiqg55crXYuXoQRKfDBFA1WEjWgP6LHhwBZeNK1VTsfTFUHCdrfp1bgwQ9xv5ski8PX9rL2dZXvgGDnw")
	if err != nil {
		t.Fatal(err)
	}

	chainPrivm0H1, err := chainPrivm0H.Child(1)
	if err != nil {
		t.Fatal(err)
	}
	chainPubm0H1 := chainPrivm0H1.Pub()

	err = testHelper(chainPrivm0H1, "xprv9wTYmMFdV23N2TdNG573QoEsfRrWKQgWeibmLntzniatZvR9BmLnvSxqu53Kw1UmYPxLgboyZQaXwTCg8MSY3H2EU4pWcQDnRnrVA1xe8fs", chainPubm0H1, "xpub6ASuArnXKPbfEwhqN6e3mwBcDTgzisQN1wXN9BJcM47sSikHjJf3UFHKkNAWbWMiGj7Wf5uMash7SyYq527Hqck2AxYysAA7xmALppuCkwQ")
	if err != nil {
		t.Fatal(err)
	}

	hKey = uint32(0x80000002)

	chainPrivm0H12H, err := chainPrivm0H1.Child(hKey)
	if err != nil {
		t.Fatal(err)
	}
	chainPubm0H12H := chainPrivm0H12H.Pub()

	err = testHelper(chainPrivm0H12H, "xprv9z4pot5VBttmtdRTWfWQmoH1taj2axGVzFqSb8C9xaxKymcFzXBDptWmT7FwuEzG3ryjH4ktypQSAewRiNMjANTtpgP4mLTj34bhnZX7UiM", chainPubm0H12H, "xpub6D4BDPcP2GT577Vvch3R8wDkScZWzQzMMUm3PWbmWvVJrZwQY4VUNgqFJPMM3No2dFDFGTsxxpG5uJh7n7epu4trkrX7x7DogT5Uv6fcLW5")
	if err != nil {
		t.Fatal(err)
	}

	chainPrivm0H12H2, err := chainPrivm0H12H.Child(2)
	if err != nil {
		t.Fatal(err)
	}
	chainPubm0H12H2 := chainPrivm0H12H2.Pub()

	err = testHelper(
		chainPrivm0H12H2,
		"xprvA2JDeKCSNNZky6uBCviVfJSKyQ1mDYahRjijr5idH2WwLsEd4Hsb2Tyh8RfQMuPh7f7RtyzTtdrbdqqsunu5Mm3wDvUAKRHSC34sJ7in334",
		chainPubm0H12H2,
		"xpub6FHa3pjLCk84BayeJxFW2SP4XRrFd1JYnxeLeU8EqN3vDfZmbqBqaGJAyiLjTAwm6ZLRQUMv1ZACTj37sR62cfN7fe5JnJ7dh8zL4fiyLHV",
	)
	if err != nil {
		t.Fatal(err)
	}

	chainPrivm0H12H210bil, err := chainPrivm0H12H2.Child(1000000000)
	if err != nil {
		t.Fatal(err)
	}
	chainPubm0H12H210bil := chainPrivm0H12H210bil.Pub()

	err = testHelper(
		chainPrivm0H12H210bil,
		"xprvA41z7zogVVwxVSgdKUHDy1SKmdb533PjDz7J6N6mV6uS3ze1ai8FHa8kmHScGpWmj4WggLyQjgPie1rFSruoUihUZREPSL39UNdE3BBDu76",
		chainPubm0H12H210bil,
		"xpub6H1LXWLaKsWFhvm6RVpEL9P4KfRZSW7abD2ttkWP3SSQvnyA8FSVqNTEcYFgJS2UaFcxupHiYkro49S8yGasTvXEYBVPamhGW6cFJodrTHy",
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestVector2Bip32(t *testing.T) {
	seed, err := hex.DecodeString("fffcf9f6f3f0edeae7e4e1dedbd8d5d2cfccc9c6c3c0bdbab7b4b1aeaba8a5a29f9c999693908d8a8784817e7b7875726f6c696663605d5a5754514e4b484542")
	if err != nil {
		t.Fatal(err)
	}

	masterprv := MasterKey(seed)
	if masterprv.String() != "xprv9s21ZrQH143K31xYSDQpPDxsXRTUcvj2iNHm5NUtrGiGG5e2DtALGdso3pGz6ssrdK4PFmM8NSpSBHNqPqm55Qn3LqFtT2emdEXVYsCzC2U" {
		t.Fatal(errors.New("master priv key doesn't match, quitting"))
	}

	masterpub := masterprv.Pub()
	if masterpub.String() != "xpub661MyMwAqRbcFW31YEwpkMuc5THy2PSt5bDMsktWQcFF8syAmRUapSCGu8ED9W6oDMSgv6Zz8idoc4a6mr8BDzTJY47LJhkJ8UB7WEGuduB" {
		t.Fatal(errors.New("master pub key doesn't match, quitting"))
	}

	chainPrivm0, err := masterprv.Child(0)
	if err != nil {
		t.Fatal(err)
	}
	chainPubm0 := chainPrivm0.Pub()

	err = testHelper(
		chainPrivm0,
		"xprv9vHkqa6EV4sPZHYqZznhT2NPtPCjKuDKGY38FBWLvgaDx45zo9WQRUT3dKYnjwih2yJD9mkrocEZXo1ex8G81dwSM1fwqWpWkeS3v86pgKt",
		chainPubm0,
		"xpub69H7F5d8KSRgmmdJg2KhpAK8SR3DjMwAdkxj3ZuxV27CprR9LgpeyGmXUbC6wb7ERfvrnKZjXoUmmDznezpbZb7ap6r1D3tgFxHmwMkQTPH",
	)
	if err != nil {
		t.Fatal(err)
	}

	var hKey uint32
	hKey = 2147483647 + 0x80000000
	chainPrivm0RndH, err := chainPrivm0.Child(hKey)
	if err != nil {
		t.Fatal(err)
	}
	chainPubm0RndH := chainPrivm0RndH.Pub()

	err = testHelper(
		chainPrivm0RndH,
		"xprv9wSp6B7kry3Vj9m1zSnLvN3xH8RdsPP1Mh7fAaR7aRLcQMKTR2vidYEeEg2mUCTAwCd6vnxVrcjfy2kRgVsFawNzmjuHc2YmYRmagcEPdU9",
		chainPubm0RndH,
		"xpub6ASAVgeehLbnwdqV6UKMHVzgqAG8Gr6riv3Fxxpj8ksbH9ebxaEyBLZ85ySDhKiLDBrQSARLq1uNRts8RuJiHjaDMBU4Zn9h8LZNnBC5y4a",
	)
	if err != nil {
		t.Fatal(err)
	}

	chainPrivm0RndH1, err := chainPrivm0RndH.Child(1)
	if err != nil {
		t.Fatal(err)
	}
	chainPubm0RndH1 := chainPrivm0RndH1.Pub()

	err = testHelper(
		chainPrivm0RndH1,
		"xprv9zFnWC6h2cLgpmSA46vutJzBcfJ8yaJGg8cX1e5StJh45BBciYTRXSd25UEPVuesF9yog62tGAQtHjXajPPdbRCHuWS6T8XA2ECKADdw4Ef",
		chainPubm0RndH1,
		"xpub6DF8uhdarytz3FWdA8TvFSvvAh8dP3283MY7p2V4SeE2wyWmG5mg5EwVvmdMVCQcoNJxGoWaU9DCWh89LojfZ537wTfunKau47EL2dhHKon",
	)
	if err != nil {
		t.Fatal(err)
	}

	hKey = 2147483646 + 0x80000000
	chainPrivm0RndH1Rnd, err := chainPrivm0RndH1.Child(hKey)
	if err != nil {
		t.Fatal(err)
	}
	chainPubm0RndH1Rnd := chainPrivm0RndH1Rnd.Pub()

	err = testHelper(
		chainPrivm0RndH1Rnd,
		"xprvA1RpRA33e1JQ7ifknakTFpgNXPmW2YvmhqLQYMmrj4xJXXWYpDPS3xz7iAxn8L39njGVyuoseXzU6rcxFLJ8HFsTjSyQbLYnMpCqE2VbFWc",
		chainPubm0RndH1Rnd,
		"xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL",
	)
	if err != nil {
		t.Fatal(err)
	}

	chainPrivm0RndH1Rnd2, err := chainPrivm0RndH1Rnd.Child(2)
	if err != nil {
		t.Fatal(err)
	}
	chainPubm0RndH1Rnd2 := chainPrivm0RndH1Rnd2.Pub()

	err = testHelper(
		chainPrivm0RndH1Rnd2,
		"xprvA2nrNbFZABcdryreWet9Ea4LvTJcGsqrMzxHx98MMrotbir7yrKCEXw7nadnHM8Dq38EGfSh6dqA9QWTyefMLEcBYJUuekgW4BYPJcr9E7j",
		chainPubm0RndH1Rnd2,
		"xpub6FnCn6nSzZAw5Tw7cgR9bi15UV96gLZhjDstkXXxvCLsUXBGXPdSnLFbdpq8p9HmGsApME5hQTZ3emM2rnY5agb9rXpVGyy3bdW6EEgAtqt",
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestVector3Bip32(t *testing.T) {
	seed, err := hex.DecodeString("4b381541583be4423346c643850da4b320e46a87ae3d2a4e6da11eba819cd4acba45d239319ac14f863b8d5ab5a0d0c64d2e8a1e7d1457df2e5a3c51c73235be")
	if err != nil {
		t.Fatal(err)
	}

	masterprv := MasterKey(seed)
	if masterprv.String() != "xprv9s21ZrQH143K25QhxbucbDDuQ4naNntJRi4KUfWT7xo4EKsHt2QJDu7KXp1A3u7Bi1j8ph3EGsZ9Xvz9dGuVrtHHs7pXeTzjuxBrCmmhgC6" {
		t.Fatal(errors.New("master priv key doesn't match, quitting"))
	}

	masterpub := masterprv.Pub()
	if masterpub.String() != "xpub661MyMwAqRbcEZVB4dScxMAdx6d4nFc9nvyvH3v4gJL378CSRZiYmhRoP7mBy6gSPSCYk6SzXPTf3ND1cZAceL7SfJ1Z3GC8vBgp2epUt13" {
		t.Fatal(errors.New("master pub key doesn't match, quitting"))
	}

	hKey := uint32(0x80000000)
	chainPrivm0H, err := masterprv.Child(hKey)
	if err != nil {
		t.Fatal(err)
	}
	chainPubm0H := chainPrivm0H.Pub()

	err = testHelper(
		chainPrivm0H,
		"xprv9uPDJpEQgRQfDcW7BkF7eTya6RPxXeJCqCJGHuCJ4GiRVLzkTXBAJMu2qaMWPrS7AANYqdq6vcBcBUdJCVVFceUvJFjaPdGZ2y9WACViL4L",
		chainPubm0H,
		"xpub68NZiKmJWnxxS6aaHmn81bvJeTESw724CRDs6HbuccFQN9Ku14VQrADWgqbhhTHBaohPX4CjNLf9fq9MYo6oDaPPLPxSb7gwQN3ih19Zm4Y",
	)
	if err != nil {
		t.Fatal(err)
	}
}
