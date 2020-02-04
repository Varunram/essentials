package xlm

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"

	utils "github.com/Varunram/essentials/utils"
	protocols "github.com/stellar/go/protocols/horizon"
)

// GetLedgerData gets the latest data from the ledger
func GetLedgerData(blockNumberx int) ([]byte, error) {
	var err error
	var data []byte

	blockNumber, err := utils.ToString(blockNumberx)
	if err != nil {
		return data, err
	}

	resp, err := http.Get(TestNetClient.HorizonURL + "ledgers/" + blockNumber)
	if err != nil || resp.Status != "200 OK" {
		return data, errors.New("API Request did not succeed")
	}

	defer func() {
		if ferr := resp.Body.Close(); ferr != nil {
			err = ferr
		}
	}()

	data, err = ioutil.ReadAll(resp.Body)
	return data, err
}

// GetBlockHash gets the block hash corresponding to a block number
func GetBlockHash(blockNumber int) (string, error) {
	var err error
	var hash string

	b, err := GetLedgerData(blockNumber)
	if err != nil {
		return hash, errors.Wrap(err, "could not get updated ledger data")
	}
	var x protocols.Ledger
	err = json.Unmarshal(b, &x)
	hash = x.Hash
	log.Printf("The block hash for block %d is: %s and the prev hash is %s", x.Sequence, hash, x.PrevHash)
	return hash, err
}

// GetLatestBlockHash gets the lastest block hash
func GetLatestBlockHash() (string, error) {
	url := TestNetClient.HorizonURL + "ledgers?cursor=now&order=desc&limit=1"
	resp, err := http.Get(url)
	if err != nil || resp.Status != "200 OK" {
		return "", errors.New("API Request did not succeed")
	}

	defer func() {
		if ferr := resp.Body.Close(); ferr != nil {
			err = ferr
		}
	}()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "could not read response body")
	}
	// hacks below follow because of stellar's incomplete go sdk support
	var x map[string]*json.RawMessage
	err = json.Unmarshal(data, &x)
	if err != nil {
		return "", errors.Wrap(err, "could not unmarshal json")
	}

	var y map[string]*json.RawMessage
	err = json.Unmarshal(*x["_embedded"], &y)
	if err != nil {
		return "", errors.Wrap(err, "could not unmarshal json")
	}

	var z []protocols.Ledger
	err = json.Unmarshal(*y["records"], &z)
	if err != nil {
		return "", errors.Wrap(err, "could not unmarshal json")
	}

	return z[0].Hash, nil
}

// GetLatestBlock fetches the latest block from the blockchain
func GetLatestBlock() (int32, error) {
	url := TestNetClient.HorizonURL + "ledgers?cursor=now&order=desc&limit=1"
	resp, err := http.Get(url)
	if err != nil || resp.Status != "200 OK" {
		return -1, errors.New("API Request did not succeed")
	}

	defer func() {
		if ferr := resp.Body.Close(); ferr != nil {
			err = ferr
		}
	}()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, errors.Wrap(err, "could not read response body")
	}
	// hacks below follow because of stellar's incomplete go sdk support
	var x map[string]*json.RawMessage
	err = json.Unmarshal(data, &x)
	if err != nil {
		return -1, errors.Wrap(err, "could not unmarshal json")
	}

	var y map[string]*json.RawMessage
	err = json.Unmarshal(*x["_embedded"], &y)
	if err != nil {
		return -1, errors.Wrap(err, "could not unmarshal json")
	}

	var z []protocols.Ledger
	err = json.Unmarshal(*y["records"], &z)
	if err != nil {
		return -1, errors.Wrap(err, "could not unmarshal json")
	}

	return z[0].Sequence, nil
}

// GetAccountData gets the account data
func GetAccountData(a string) ([]byte, error) {
	var err error
	var data []byte
	resp, err := http.Get(TestNetClient.HorizonURL + "accounts/" + a)
	if err != nil {
		return data, errors.Wrap(err, "could not get /accounts/ endpoint from API")
	}
	if resp.Status != "200 OK" {
		return data, errors.New("Request did not succeed")
	}

	defer func() {
		if ferr := resp.Body.Close(); ferr != nil {
			err = ferr
		}
	}()

	data, err = ioutil.ReadAll(resp.Body)
	return data, err
}

// GetNativeBalance gets the xlm balance of a specific account
func GetNativeBalance(publicKey string) float64 {
	var err error
	b, err := GetAccountData(publicKey)
	if err != nil {
		log.Println(err)
		return -1
	}
	var x protocols.Account
	err = json.Unmarshal(b, &x)
	if err != nil {
		log.Println(err)
		return -1
	}
	for _, balance := range x.Balances {
		if balance.Asset.Type == "native" {
			balance, err := utils.ToFloat(balance.Balance)
			if err != nil {
				return 0
			}
			return balance
		}
	}

	log.Println("Native balance not found")
	return 0
}

// GetAssetBalance gets the balance of the user in the specific asset
func GetAssetBalance(publicKey string, assetName string) float64 {
	var err error
	b, err := GetAccountData(publicKey)
	if err != nil {
		log.Println(err)
		return -1
	}
	var x protocols.Account
	err = json.Unmarshal(b, &x)
	if err != nil {
		log.Println(err)
		return -1
	}
	for _, balance := range x.Balances {
		if balance.Asset.Code == assetName {
			balance, err := utils.ToFloat(balance.Balance)
			if err != nil {
				return 0
			}
			return balance
		}
	}

	// if the asset does not exist, the balance associated is zero
	// but to differentiate between not having the asset and never
	// having had the asset, we return -1
	log.Println("Asset balance not found")
	return 0
}

// GetAssetTrustLimit gets the trust limit that the user has with an issuer
func GetAssetTrustLimit(publicKey string, assetName string) float64 {
	var err error
	b, err := GetAccountData(publicKey)
	if err != nil {
		return -1
	}
	var x protocols.Account
	err = json.Unmarshal(b, &x)
	if err != nil {
		return -1
	}
	for _, balance := range x.Balances {
		if balance.Asset.Code == assetName {
			balance, err := utils.ToFloat(balance.Balance)
			if err != nil {
				return 0
			}
			return balance
		}
	}
	return 0
}

// GetAllBalances gets all the balances associated with a certain account.
func GetAllBalances(publicKey string) ([]protocols.Balance, error) {
	account, err := ReturnSourceAccountPubkey(publicKey)
	if err != nil {
		log.Println(err)
		return nil, errors.Wrap(err, "could not load account")
	}
	return account.Balances, nil
}

// HasStableCoin checks whether the PublicKey has a stablecoin balance
func HasStableCoin(publicKey string) bool {
	account, err := ReturnSourceAccountPubkey(publicKey)
	if err != nil {
		return false
	}

	for _, balance := range account.Balances {
		if balance.Asset.Code == "STABLEUSD" {
			return true
		}
	}
	return false
}

// GetTransactionData gets tx data
func GetTransactionData(txhash string) ([]byte, error) {
	var err error
	var data []byte
	resp, err := http.Get(TestNetClient.HorizonURL + "transactions/" + txhash)
	if err != nil || resp.Status != "200 OK" {
		// check here since if we don't, we need to check the body of the unmarshalled
		// response to see if we have 0
		return data, errors.New("API Request did not succeed")
	}

	defer func() {
		if ferr := resp.Body.Close(); ferr != nil {
			err = ferr
		}
	}()

	data, err = ioutil.ReadAll(resp.Body)
	return data, err
}

// GetTransactionHeight gets height at which a tx was confirmed
func GetTransactionHeight(txhash string) (int, error) {
	var err error
	var txheight int

	b, err := GetTransactionData(txhash)
	if err != nil {
		return txheight, errors.Wrap(err, "could not get transaction data")
	}
	var x protocols.Transaction
	err = json.Unmarshal(b, &x)
	return int(x.Ledger), err
}
