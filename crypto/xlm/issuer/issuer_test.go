// +build all travis

package issuer

import (
	"os"
	"testing"

	xlm "github.com/Varunram/essentials/crypto/xlm"
)

func TestIssuer(t *testing.T) {
	var err error
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	TestDir := wd + "/test"
	os.MkdirAll(TestDir, 0775)
	platformSeed, platformPubkey, err := xlm.GetKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	err = xlm.GetXLM(platformPubkey)
	if err != nil {
		t.Fatal(err)
	}
	err = InitIssuer(TestDir, 1, "blah")
	if err != nil {
		t.Fatal(err)
	}
	err = FundIssuer(TestDir, 1, "blah", platformSeed)
	if err != nil {
		t.Fatal(err)
	}
	err = FundIssuer(TestDir, 1, "cool", platformSeed)
	if err == nil {
		t.Fatalf("not able to catch invalid seed error, quitting!")
	}
	_, err = FreezeIssuer(TestDir, 1, "blah")
	if err != nil {
		t.Fatal(err)
	}
	err = DeleteIssuer(TestDir, 1)
	if err != nil {
		t.Fatal(err)
	}
	os.RemoveAll(TestDir)
}
