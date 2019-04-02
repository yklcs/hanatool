package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Client brings in the HTTP Client from reserver.go
var Client *http.Client

func Login(id, pw string) []byte {

	url := "http://gotest.hana.hs.kr:8081/json/loginProc.ajax"

	payloadRaw := fmt.Sprintf("mUsr_ID=%s&mUsr_PW=%s&loginArr=&push_Token=&undefined=", id, pw)
	payload := strings.NewReader(payloadRaw)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, _ := Client.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}

func Response(byt []byte) bool {
	var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat["result"])
	return true
}
