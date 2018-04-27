package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	gorilla "github.com/vharish836/rpc/json"
)

type getinforesp struct {
	Version         string  `json:"version"`
	Nodeversion     int     `json:"nodeversion"`
	Protocolversion int     `json:"protocolversion"`
	Chainname       string  `json:"chainname"`
	Description     string  `json:"description"`
	Protocol        string  `json:"protocol"`
	Port            int     `json:"port"`
	Setupblocks     int     `json:"setupblocks"`
	Nodeaddress     string  `json:"nodeaddress"`
	Burnaddress     string  `json:"burnaddress"`
	Incomingpaused  bool    `json:"incomingpaused"`
	Miningpaused    bool    `json:"miningpaused"`
	Walletversion   int     `json:"walletversion"`
	Balance         float64 `json:"balance"`
	Walletdbversion int     `json:"walletdbversion"`
	Reindex         bool    `json:"reindex"`
	Blocks          int     `json:"blocks"`
	Timeoffset      int     `json:"timeoffset"`
	Connections     int     `json:"connections"`
	Proxy           string  `json:"proxy"`
	Difficulty      float64 `json:"difficulty"`
	Testnet         bool    `json:"testnet"`
	Keypoololdest   int     `json:"keypoololdest"`
	Keypoolsize     int     `json:"keypoolsize"`
	Paytxfee        float64 `json:"paytxfee"`
	Relayfee        float64 `json:"relayfee"`
	Errors          string  `json:"errors"`
}

type GetRuntimeParamReq []string

func main() {
	url := "http://localhost:8383/"
	params := &GetRuntimeParamReq{"gen", "false"}
	jbuf, err := gorilla.EncodeClientRequest("BridgeService.SetRuntimeParam",params)
	if err != nil {
		log.Fatalf("could not encode. %s", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jbuf))
	if err != nil {
		log.Fatalf("Could not create new request. %s", err)
	}
	req.SetBasicAuth("username", "password")
	req.Header.Set("Content-Type", "application/json")
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("sending request failed. %s", err)
	}
	defer rsp.Body.Close()
	var result getinforesp
	err = gorilla.DecodeClientResponse(rsp.Body, &result)
	if err != nil {
		log.Fatalf("could not decode response. %s", err)
	}
	jbuf, err = json.MarshalIndent(result, "", "\t")
	if err != nil {
		log.Fatalf("could not encode request. %s", err)
	}
	log.Printf("buffer: %s", jbuf)
}
