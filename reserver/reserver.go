/*
TODO: Edit API to use go.hana.hs.kr




*/
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"sync"
	"time"
)

var cookieJar, _ = cookiejar.New(nil)
var client = &http.Client{Jar: cookieJar}

func main() {
	fmt.Printf("Program starting at %s\n", time.Now().Format(time.RFC3339Nano))

	tries := 25
	var wg1, wg2 sync.WaitGroup

	id := flag.String("id", "", "Username/ID")
	pw := flag.String("pw", "", "Password")
	sCode := flag.Int("seat", 84, "seat code")
	tCode1 := flag.Int("time1", 1, "time code 1")
	tCode2 := flag.Int("time2", 4, "time code 2")
	min := flag.Int("min", 59, "minute")
	sec := flag.Int("sec", 59, "second")

	flag.Parse()

	if *id == "" || *pw == "" {
		fmt.Println("Please re-check all flags.")
		os.Exit(1)
	}

	// tt := time.Now()
	fmt.Println("Logging in...")
	login(*id, *pw)
	fmt.Printf("Logged in and beginning wait at %s\n", time.Now().Format(time.RFC3339Nano))

	for i := 0; i < 2000; i++ {
		if time.Now().Minute() == *min && time.Now().Second() == *sec {
			t := time.Now()
			for j := 0; j < tries; j++ {
				wg1.Add(1)
				wg2.Add(1)
				go func() {
					defer wg1.Done()
					reserve(*sCode, *tCode1)
				}()
				go func() {
					defer wg2.Done()
					reserve(*sCode, *tCode2)
				}()
				time.Sleep(50 * time.Millisecond)
			}
			fmt.Printf("\n%s\n", time.Since(t))

			wg1.Wait()
			wg2.Wait()

			break
		} else {
			time.Sleep(50 * time.Millisecond)
		}
	}

	fmt.Printf("Exiting program at %s", time.Now().Format(time.RFC3339Nano))
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

func reserve(sCode, tCode int) {
	url := "http://hi.hana.hs.kr/SYSTEM_Plan/Lib_System/Lib_System_Reservation/reservation_exec.asp"

	rawPayload := fmt.Sprintf("code=001&s_code=%d&t_code=%d", sCode, tCode)
	payload := strings.NewReader(rawPayload)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)

	if err != nil {
		fmt.Printf("Something went wrong...")
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))

	if strings.Contains(string(body), "되었습니다") {
		fmt.Printf("Success at %s for seat %d, time %d\n", time.Now().Format(time.RFC3339Nano), sCode, tCode)
	} else {
		fmt.Printf("Error at %s for seat %d, time %d\n", time.Now().Format(time.RFC3339Nano), sCode, tCode)
	}
	return
}
