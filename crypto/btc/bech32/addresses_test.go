package bech32

import (
	"encoding/hex"
	"log"
	"testing"
)

func TestAddr(t *testing.T) {
	var bech32Adr = []string{"tb1qmsu7ck0tun9qe2wgthu35xcu6asa5aq5tejh02", "bc1q6sh5tzw0c650hutmm58s7srdut8qrg05a4kfmd"}
	var base58Adr = []string{"n1bQLqMS86rjrWLVkN78FdaCzjMjQwZ2k1", "1LLvtAD6EbNcSc5QXMo46bAdWkEHmbM8xg"}

	for i, _ := range bech32Adr {
		base58Conv, err := Bech32ToBase58Addr(bech32Adr[i][0:2], bech32Adr[i])
		if err != nil {
			log.Fatal(err)
		}

		if base58Conv != base58Adr[i] {
			t.Fatalf("converted base 58 did not match with required base58 string")
		}

		bech32Conv, err := Base58ToBech32Address(base58Conv)
		if err != nil {
			log.Fatal(err)
		}

		if bech32Conv != bech32Adr[i] {
			t.Fatalf("converted bech32 address does not match with required bech32 address")
		}
	}

	privkey, err := hex.DecodeString("0c28fca386c7a227600b2fe50b7cae11ec86d3bf1fbe471be89827e19d72aa1d")
	if err != nil {
		t.Fatal(err)
	}
	wif, err := PrivKeyToWIF("mainnet", false, privkey)
	if err != nil {
		t.Fatal(err)
	}
	if wif != "5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ" {
		t.Fatalf("wif doesn't match with expected value")
	}

	privKeyCheckByte, err := WIFToPrivateKey(wif, false)
	if err != nil {
		t.Fatal(err)
	}

	privKeyCheck := hex.EncodeToString(privKeyCheckByte)
	if privKeyCheck != "0c28fca386c7a227600b2fe50b7cae11ec86d3bf1fbe471be89827e19d72aa1d" {
		log.Println(privKeyCheck)
		t.Fatalf("decoded privkey not equal to one obtained from wif")
	}

	if CheckCheckSum(wif) != nil {
		t.Fatalf("checksums don't match, quitting")
	}
}
