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

func main() {
	data, err := GetChainTxStats("user", "password", "1")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(data))
}
