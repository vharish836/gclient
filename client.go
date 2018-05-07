package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

// JSONRequest ...
type JSONRequest struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
	ID     interface{}   `json:"id"`
}

// JSONResponse ...
type JSONResponse struct {
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
	ID     interface{} `json:"id"`
}

func printmap(m map[string]interface{}) {
	rbuf, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		log.Fatalf("could not encode request. %s", err)
	}
	fmt.Printf("\n%s\n", rbuf)
}

func parseparam(s string) interface{} {
	var result interface{}
	var err error
	result, err = strconv.ParseInt(s, 0, 0)
	if err == nil {
		return result
	}
	result, err = strconv.ParseBool(s)
	if err == nil {
		return result
	}
	result, err = strconv.ParseFloat(s, 64)
	if err == nil {
		return result
	}
	return s
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
	req := JSONRequest{Method: params[0], Params: make([]interface{}, 0), ID: rand.Int()}

	for i := 1; i < len(params); i++ {
		req.Params = append(req.Params, parseparam(params[i]))
	}

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
			log.Fatalf("coudl not decode response: %s", err)
		}
		if result.Result != nil {
			fmt.Printf("Response (result) ==>")
			switch result.Result.(type) {
			case map[string]interface{}:
				printmap(result.Result.(map[string]interface{}))
			case []interface{}:
				for _, v := range result.Result.([]interface{}) {
					switch v.(type) {
					case map[string]interface{}:
						printmap(v.(map[string]interface{}))
					case string:
						fmt.Printf("%s", v.(string))
					}
				}
			case string:
				fmt.Printf("%s", result.Result.(string))
			}
		}
		if result.Error != nil {
			fmt.Printf("Response (error) ==>")
			switch result.Error.(type) {
			case map[string]interface{}:
				printmap(result.Error.(map[string]interface{}))
			case []interface{}:
				for _, v := range result.Error.([]interface{}) {
					switch v.(type) {
					case map[string]interface{}:
						printmap(v.(map[string]interface{}))
					case string:
						fmt.Printf("%s", v.(string))
					}
				}
			case string:
				fmt.Printf("%s", result.Error.(string))
			}
		}
	}
}
