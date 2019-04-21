/*
TODO: Edit API to use go.hana.hs.kr




*/
package main

import (
	"net/http"
	"net/http/cookiejar"

	"github.com/rocketll/hanatool/utils"
)

var cookieJar, _ = cookiejar.New(nil)
var client = &http.Client{Jar: cookieJar}

func main() {
	utils.Client = client
	res := utils.Login("id", "pw")
	utils.Response(res)

}
