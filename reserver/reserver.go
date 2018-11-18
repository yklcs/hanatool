package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"time"
)

var cookieJar, _ = cookiejar.New(nil)

var client = &http.Client{Jar: cookieJar}

func main() {
	id := flag.String("id", "", "Username/ID")
	pw := flag.String("pw", "", "Password")
	sCode := flag.Int("seat", 96, "seat code")
	tCode1 := flag.Int("time1", 7, "time code 1")
	tCode2 := flag.Int("time2", 9, "time code 2")
	debug := flag.Bool("debug", false, "debugging mode")

	flag.Parse()

	if *id == "" || *pw == "" {
		fmt.Println("Please re-check all flags.")
		os.Exit(1)
	}

	if *debug {
		fmt.Println("Beginning delay.")
		t := time.Now()
		time.Sleep(56 * time.Second)
		fmt.Printf("%s delay done.\n", time.Since(t))
	}

	tt := time.Now()
	fmt.Println("Logging in...")
	login(*id, *pw)
	fmt.Printf("%s to attempt login.\n", time.Since(tt))

	ok := false
	n := 1

	for !ok && n <= 25 {
		ttt := time.Now()
		if reserve(*sCode, *tCode1) && reserve(*sCode, *tCode2) {
			ok = true
			break
		}
		fmt.Printf("Request took %s\n", time.Since(ttt))
		fmt.Printf("Failed, retrying... Attempt %d\n", n)
		n++
		time.Sleep(200 * time.Millisecond)
	}

	if ok {
		fmt.Println("Success!")
		os.Exit(0)
	} else {
		fmt.Printf("Error... Exited.")
	}
}

func login(id, pw string) *http.Response {

	url := "http://hi.hana.hs.kr/proc/login_proc.asp"

	payloadRaw := fmt.Sprintf("login_id=%s&login_pw=%s&x=0&y=0", id, pw)
	payload := strings.NewReader(payloadRaw)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, _ := client.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	if strings.Contains(string(body), "게시판") {
		fmt.Println("Logged in!")
	} else {
		fmt.Println("Couldn't log in")
	}
	return res
}

func reserve(sCode, tCode int) bool {
	url := "http://hi.hana.hs.kr/SYSTEM_Plan/Lib_System/Lib_System_Reservation/reservation_exec.asp"

	rawPayload := fmt.Sprintf("code=001&s_code=%d&t_code=%d", sCode, tCode)
	payload := strings.NewReader(rawPayload)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)

	if err != nil {
		fmt.Println("Something went wrong...")
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	if strings.Contains(string(body), "되었습니다") {
		fmt.Println("Successful.")
		return true
	} else {
		fmt.Printf("Error... ")
		return false
	}
}