package assets

import (
	//"log"
	"github.com/pkg/errors"

	xlm "github.com/Varunram/essentials/crypto/xlm"
	utils "github.com/Varunram/essentials/utils"
	"github.com/stellar/go/network"
	build "github.com/stellar/go/txnbuild"
)

// AssetID generates a new stellar compatible asset
func AssetID(inputString string) string {
	// so the assetID right now is a hash of the asset name, concatenated investor public keys and nonces
	x := utils.SHA3hash(inputString)
	return "OXA" + x[64:73] // max length of an asset in stellar is 12 (OXA: OpenX Asset)
}

// CreateAsset creates a new asset belonging to the passed public key
func CreateAsset(assetName string, PublicKey string) build.Asset {
	// need to set a couple flags here
	return build.CreditAsset{assetName, PublicKey}
}

// TrustAsset trusts an asset issued by an account and signs a transaction with a
// preset limit on how much it is willing to trust the issuer
func TrustAsset(assetCode string, assetIssuer string, limitI float64, seed string) (string, error) {
	// TRUST is FROM Seed TO assetIssuer
	passphrase := network.TestNetworkPassphrase
	sourceAccount, mykp, err := xlm.ReturnSourceAccount(seed)
	if err != nil {
		return "", err
	}

	/*
	var limit string
	switch limitInput.(type) {
	case string:
		limit = limitInput.(string)
	case int:
		limit, err = utils.ToString(limitInput.(int))
		if err != nil {
			return "", errors.New("could not convert limit to string")
		}
	}
	*/
	limit, err := utils.ToString(limitI)
	if err != nil {
		return "", errors.New("could not convert limit to string")
	}

	op := build.ChangeTrust{
		Line:  build.CreditAsset{assetCode, assetIssuer},
		Limit: limit,
	}

	tx := build.Transaction{
		SourceAccount: &sourceAccount,
		Operations:    []build.Operation{&op},
		Timebounds:    build.NewInfiniteTimeout(),
		Network:       passphrase,
	}

	_, txHash, err := xlm.SendTx(mykp, tx)
	if err != nil {
		return "", err
	}

	return txHash, err
}

// SendAssetFromIssuer transfers an asset from the issuer to the desired publickey.
func SendAssetFromIssuer(assetCode string, destination string, amount string,
	seed string, issuerPubkey string) (int32, string, error) {

	passphrase := network.TestNetworkPassphrase
	sourceAccount, mykp, err := xlm.ReturnSourceAccount(seed)
	if err != nil {
		return -1, "", err
	}

	op := build.Payment{
		Destination: destination,
		Amount:      amount,
		Asset:       build.CreditAsset{assetCode, issuerPubkey},
	}

	tx := build.Transaction{
		SourceAccount: &sourceAccount,
		Operations:    []build.Operation{&op},
		Timebounds:    build.NewInfiniteTimeout(),
		Network:       passphrase,
	}

	return xlm.SendTx(mykp, tx)
}

// SendAssetToIssuer sends an asset back to the issuer
func SendAssetToIssuer(assetCode string, destination string, amount string,
	seed string) (int32, string, error) {

	passphrase := network.TestNetworkPassphrase
	sourceAccount, mykp, err := xlm.ReturnSourceAccount(seed)
	if err != nil {
		return -1, "", err
	}

	op := build.Payment{
		Destination: destination,
		Amount:      amount,
		Asset:       build.CreditAsset{assetCode, destination},
	}

	tx := build.Transaction{
		SourceAccount: &sourceAccount,
		Operations:    []build.Operation{&op},
		Timebounds:    build.NewInfiniteTimeout(),
		Network:       passphrase,
	}

	return xlm.SendTx(mykp, tx)
}

// SendAsset sends an asset to a destination which has an established trustline with the issuer
func SendAsset(assetCode string, issuerPubkey string, destination string, amount string,
	seed string, memo string) (int32, string, error) {
	passphrase := network.TestNetworkPassphrase
	sourceAccount, mykp, err := xlm.ReturnSourceAccount(seed)
	if err != nil {
		return -1, "", err
	}

	op := build.Payment{
		Destination: destination,
		Amount:      amount,
		Asset:       build.CreditAsset{assetCode, issuerPubkey},
	}

	tx := build.Transaction{
		SourceAccount: &sourceAccount,
		Operations:    []build.Operation{&op},
		Timebounds:    build.NewInfiniteTimeout(),
		Network:       passphrase,
		Memo:          build.Memo(build.MemoText(memo)),
	}

	return xlm.SendTx(mykp, tx)
}
