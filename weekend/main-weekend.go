package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

var cookieJar, _ = cookiejar.New(nil)

var client = &http.Client{Jar: cookieJar}

func main() {
	id := flag.String("id", "", "Username/ID")
	pw := flag.String("pw", "", "Password")
	// isWeekend := flag.Bool("weekend", false, "Is weekend")
	// sCode := flag.String("seat", "", "Seat")
	// tCode := flag.Int("time", 1, "time")

	flag.Parse()
	// 0 -3

	tt := time.Now()

	time.Sleep(58 * time.Second)
	// 58 55

	login(*id, *pw)

	ok := false
	n := 0

	for !ok && n < 20 {
		t := time.Now()
		if reserve(96, 7) && reserve(96, 9) {
			ok = true
			break
		}
		fmt.Printf("Failed, retrying... Attempt %d\n", n)
		n++
		time.Sleep(200 * time.Millisecond)
		fmt.Println(time.Since(t))
	}

	fmt.Println(time.Since(tt))
	if ok {
		fmt.Println("Success!")
	} else {
		fmt.Printf("Error...")
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
	fmt.Println(res)
	// fmt.Println(string(body))
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

	fmt.Println(res)
	fmt.Println(string(body))

	if strings.Contains(string(body), "되었습니다") {
		fmt.Println("Successful.")
		return true
	} else {
		fmt.Printf("Error...")
		return false
	}
}
