/*
TODO: Edit API to use go.hana.hs.kr




*/
package main

import (
	"net/http"
	"net/http/cookiejar"

	"hanatool/utils"
)

var cookieJar, _ = cookiejar.New(nil)
var client = &http.Client{Jar: cookieJar}

func main() {
	utils.Client = client
	res := utils.Login("rocketll", "yk8525yk")
	utils.Response(res)

}