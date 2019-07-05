package main

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"

	utils "github.com/Varunram/essentials/utils"
)

type RPCReq struct {
	JsonRPC string      `json:"jsonrpc"`
	ID      string      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

var BitcoindURL = "http://localhost:18443/" // for regtest

func PostReq(username string, password string, payload RPCReq) ([]byte, error) {
	var req *http.Request
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal json, quitting")
	}

	req, err = http.NewRequest("POST", BitcoindURL, bytes.NewBuffer(payloadJson))
	if err != nil {
		return nil, errors.Wrap(err, "did not POST to bitcoind")
	}

	req.SetBasicAuth(username, password)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "did not make http request to bitcoind")
	}

	defer res.Body.Close()
	x, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "did not read from ioutil")
	}

	return x, nil
}

func GetBestBlockHash(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getbestblockhash"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func GetBlock(username string, password string, blockhash string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getblock"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{blockhash}

	return PostReq(username, password, payload)
}

func GetBlockCount(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getblockcount"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func GetBlockHash(username string, password string, blockNumber uint32) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getblockhash"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []uint32{blockNumber}

	return PostReq(username, password, payload)
}

func GetBlockHeader(username string, password string, blockhash string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getblockheader"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{blockhash}

	return PostReq(username, password, payload)
}

func GetBlockStats(username string, password string, hashOrHeight string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getblockstats"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	height, err := utils.StoICheck(hashOrHeight)
	if err != nil {
		payload.Params = []string{hashOrHeight}
	} else {
		payload.Params = []int{height}
	}

	return PostReq(username, password, payload)
}

func GetChainTips(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getchaintips"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func GetChainTxStats(username string, password string, nBlocks string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getchaintxstats"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	nBlocksInt, err := utils.StoICheck(nBlocks)
	if err != nil {
		return nil, errors.New("input block height not integer")
	}
	payload.Params = []int{nBlocksInt}

	return PostReq(username, password, payload)
}

func GetDifficulty(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getdifficulty"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func GetMempoolAncestors(username string, password string, txid string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getmempoolancestors"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{txid}

	return PostReq(username, password, payload)
}

func GetMempoolEntry(username string, password string, txid string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getmempoolentry"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{txid}

	return PostReq(username, password, payload)
}

func GetMempoolInfo(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getmempoolinfo"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func GetRawMempool(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getrawmempool"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func GetTxOut(username string, password string, txid string, n int) ([]byte, error) {
	var payload RPCReq
	payload.Method = "gettxout"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = [2]interface{}{txid, n}

	return PostReq(username, password, payload)
}

// TODO: fix this route
func GetTxOutProof(username string, password string, txid string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "gettxoutproof"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	temp := make(map[string]interface{})
	temp["txids"] = txid
	payload.Params = temp

	return PostReq(username, password, payload)
}

func GetTxOutSetInfo(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "gettxoutsetinfo"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func PreciousBlock(username string, password string, blockhash string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "preciousblock"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{blockhash}

	return PostReq(username, password, payload)
}

func PruneBlockchain(username string, password string, height string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getchaintxstats"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	heightI, err := utils.StoICheck(height)
	if err != nil {
		return nil, errors.New("input block height not integer")
	}
	payload.Params = []int{heightI}

	return PostReq(username, password, payload)
}

func SaveMempool(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "SaveMempool"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

// TODO: implement this route
func ScanTxOutset(username string, password string) ([]byte, error) {
	return nil, nil
}

func VerifyChain(username string, password string, nBlocks string, checkLevel string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "verifychain"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	nBlocksInt, err := utils.StoICheck(nBlocks)
	if err == nil {
		payload.Params = []int{nBlocksInt}
	}

	checkLevelInt, err := utils.StoICheck(checkLevel)
	if err == nil {
		payload.Params = []int{checkLevelInt}
	}

	return PostReq(username, password, payload)
}

// TODO: implement this route
func VerifyTxOutProof(username string, password string, txproof string) ([]byte, error) {
	return nil, nil
}

func main() {
	data, err := VerifyChain("user", "password", "1", "1")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(data))
}
