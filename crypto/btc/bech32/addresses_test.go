package bech32

import (
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
}
