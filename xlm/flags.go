package xlm

import (
	build "github.com/stellar/go/txnbuild"
)

// SetAuthImmutable sets the auth_immutable flag on an account
func SetAuthImmutable(seed string) (int32, string, error) {
	//  Create with Auth immutable since we don't want the asset to be revocable
	sourceAccount, mykp, err := ReturnSourceAccount(seed)
	if err != nil {
		return -1, "", err
	}

	op := build.SetOptions{
		SetFlags: []build.AccountFlag{build.AuthImmutable},
	}

	memo := "set immutable"

	return SendTx(mykp, &sourceAccount, memo, build.Operation(&op))
}

// FreezeAccount freezes an account
func FreezeAccount(seed string) (int32, string, error) {
	sourceAccount, mykp, err := ReturnSourceAccount(seed)
	if err != nil {
		return -1, "", err
	}

	op := build.SetOptions{
		MasterWeight:    build.NewThreshold(0),
		LowThreshold:    build.NewThreshold(0),
		MediumThreshold: build.NewThreshold(0),
		HighThreshold:   build.NewThreshold(0),
	}

	memo := "freeze account"

	return SendTx(mykp, &sourceAccount, memo, build.Operation(&op))
}
