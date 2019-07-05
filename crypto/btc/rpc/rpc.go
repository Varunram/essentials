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
	payload.Method = "savemempool"
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

func GetMemoryInfo(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getmemoryinfo"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func GetRPCInfo(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getrpcinfo"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func Help(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "help"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func Logging(username string, password string, params ...string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "logging"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	var arr []string
	for _, param := range params {
		arr = append(arr, param)
	}
	payload.Params = arr

	return PostReq(username, password, payload)
}

func Stop(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "stop"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func Uptime(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "uptime"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func Generate(username string, password string, nBlocks string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "generate"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	nBlocksInt, err := utils.StoICheck(nBlocks)
	if err == nil {
		payload.Params = []int{nBlocksInt}
	}

	return PostReq(username, password, payload)
}

func GenerateToAddress(username string, password string, nBlocks string, address string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "generatetoaddress"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	nBlocksInt, err := utils.StoICheck(nBlocks)
	if err == nil {
		payload.Params = [2]interface{}{nBlocksInt, address}
	}

	return PostReq(username, password, payload)
}

func SubmitBlock(username string, password string, hexdata string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "submitblock"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{hexdata}

	return PostReq(username, password, payload)
}

func SubmitHeader(username string, password string, hexdata string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "submitheader"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{hexdata}

	return PostReq(username, password, payload)
}

func AddNode(username string, password string, node string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "addnode"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{node}

	return PostReq(username, password, payload)
}

func ClearBanned(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "clearbanned"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func DisconnectNode(username string, password string, address string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "disconnectnode"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{address}

	return PostReq(username, password, payload)
}

func GetAddedNodeInfo(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getaddednodeinfo"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func GetConnectionCount(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getconnectioncount"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func GetNetTotals(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getnettotals"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func GetNetworkInfo(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getnetworkinfo"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func GetNodeAddresses(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getnodeaddresses"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func GetPeerInfo(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getpeerinfo"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func ListBanned(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "listbanned"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func Ping(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "ping"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

// TODO: SetBan

func SetNetworkActive(username string, password string, state bool) ([]byte, error) {
	var payload RPCReq
	payload.Method = "setnetworkactive"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []interface{}{state}

	return PostReq(username, password, payload)
}

func main() {
	data, err := SetNetworkActive("user", "password", true)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(data))
}
