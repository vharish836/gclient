package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

// JSONRequest ...
type JSONRequest struct {
	Method string      `json:"method"`
	Params []string    `json:"params"`
	ID     interface{} `json:"id"`
}

// JSONResponse ...
type JSONResponse struct {
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
	ID     interface{} `json:"id"`
}

func main() {
	addr := flag.String("addr", "http://localhost:8383", "server address")
	username := flag.String("username", "username", "User name")
	password := flag.String("password", "password", "Password")

	flag.Parse()

	params := flag.Args()
	if len(params) == 0 {
		log.Fatalf("please provide method and its args\n")
	}
	var err error
	var jbuf []byte
	req := JSONRequest{Method: params[0], Params: params[1:], ID: rand.Int()}
	jbuf, err = json.Marshal(&req)
	if err != nil {
		log.Fatalf("could not encode. %s", err)
	}
	fmt.Printf("Request <==\n%s\n", jbuf)
	hreq, err := http.NewRequest("POST", *addr, bytes.NewBuffer(jbuf))
	if err != nil {
		log.Fatalf("Could not create new request. %s", err)
	}
	hreq.SetBasicAuth(*username, *password)
	hreq.Header.Set("Content-Type", "application/json")
	rsp, err := http.DefaultClient.Do(hreq)
	if err != nil {
		log.Fatalf("sending request failed. %s", err)
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != 200 {
		rbuf := make([]byte, 200)
		n, _ := rsp.Body.Read(rbuf)
		if n == 0 {
			fmt.Printf("Response (Status)==>\n%s\n", rsp.Status)
		} else {
			fmt.Printf("Response ==>\n%s\n", rbuf)
		}
	} else {
		result := JSONResponse{}
		err = json.NewDecoder(rsp.Body).Decode(&result)
		if err != nil {
			fmt.Printf("Response(error) ==>\n%s\n", err)
		} else {
			robj, ok := result.Result.(map[string]interface{})
			if ok == true {
				rbuf, err := json.MarshalIndent(robj, "", "\t")
				if err != nil {
					log.Fatalf("could not encode request. %s", err)
				}
				fmt.Printf("Response (Result) ==>\n%s\n", rbuf)
			} else {
				robj, ok := result.Error.(map[string]interface{})
				if ok == true {
					rbuf, err := json.MarshalIndent(robj, "", "\t")
					if err != nil {
						log.Fatalf("could not encode request. %s", err)
					}
					fmt.Printf("Response (Error) ==>\n%s\n", rbuf)
				} else {
					if result.Result != nil {
						fmt.Printf("Response (Result) ==>\n%s\n", result.Result)
					} else if result.Error != nil {
						fmt.Printf("Response (Error) ==>\n%s\n", result.Error)
					} else {
						fmt.Printf("Response () ==> \n%s\n", "empty")
					}					
				}				
			}
		}
	}
}
