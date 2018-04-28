package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	gorilla "github.com/vharish836/rpc/json"
)

func main() {
	addr := flag.String("addr", "http://localhost:8383", "server address")
	username := flag.String("username", "username", "User name")
	password := flag.String("password", "password", "Password")
	method := flag.String("method", "BridgeService.GetInfo", "Method to invoke followed by its args")

	flag.Parse()

	args := flag.Args()
	var params []interface{}
	var err error
	var jbuf []byte
	for i := 0; i < len(args); i++ {
		params = append(params, args[i])
	}
	if args == nil {
		jbuf, err = gorilla.EncodeClientRequest(*method, nil)
	} else {
		jbuf, err = gorilla.EncodeClientRequest(*method, params)
	}

	if err != nil {
		log.Fatalf("could not encode. %s", err)
	}
	fmt.Printf("Request <==\n%s\n", jbuf)
	req, err := http.NewRequest("POST", *addr, bytes.NewBuffer(jbuf))
	if err != nil {
		log.Fatalf("Could not create new request. %s", err)
	}
	req.SetBasicAuth(*username, *password)
	req.Header.Set("Content-Type", "application/json")
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("sending request failed. %s", err)
	}
	defer rsp.Body.Close()
	var result interface{}
	err = gorilla.DecodeClientResponse(rsp.Body, &result)
	if err != nil {
		fmt.Printf("Response ==>\n%s\n", err)
	} else {
		jbuf, err = json.MarshalIndent(result, "", "\t")
		if err != nil {
			log.Fatalf("could not encode request. %s", err)
		}
		fmt.Printf("Response ==>\n%s\n", jbuf)
	}
}
