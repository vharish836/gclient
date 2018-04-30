package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	gorilla "github.com/gorilla/rpc/json"
)

func main() {
	addr := flag.String("addr", "http://localhost:8383", "server address")
	username := flag.String("username", "username", "User name")
	password := flag.String("password", "password", "Password")

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Fatalf("please provide method and its args\n")
	}
	var params []interface{}
	var err error
	var jbuf []byte
	for i := 1; i < len(args); i++ {
		params = append(params, args[i])
	}
	method := args[0]
	if len(args) == 1 {
		jbuf, err = gorilla.EncodeClientRequest(method, nil)
	} else {
		jbuf, err = gorilla.EncodeClientRequest(method, params)
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
	if rsp.StatusCode != 200 {
		rbuf := make([]byte,200)
		n,_ := rsp.Body.Read(rbuf)
		if n == 0 {
			fmt.Printf("Response (Status)==>\n%s\n", rsp.Status)
		} else {
			fmt.Printf("Response ==>\n%s\n",rbuf)
		}		
	} else {
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
}
