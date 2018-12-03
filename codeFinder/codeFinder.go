package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

var cookieJar, _ = cookiejar.New(nil)

var client = &http.Client{Jar: cookieJar}

func main() {
	id := flag.String("id", "", "Username/ID")
	pw := flag.String("pw", "", "Password")
	isWeekend := flag.Bool("weekend", false, "Is weekend")
	seat := flag.String("seat", "", "Seat")
	time := flag.Int("time", 1, "time")

	flag.Parse()

	if *id == "" || *pw == "" || *seat == "" {
		fmt.Println("Please re-check all flags.")
		os.Exit(1)
	}

	login(*id, *pw)

	tCode := getTimeCode(*time, *isWeekend)
	sCode := getSeatCode(*seat)

	fmt.Printf("Seat code is %d, Time code is %d", sCode, tCode)
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

func getTimeCode(time int, isWeekEnd bool) int {
	var tCode int
	if isWeekEnd {
		if time == 1 {
			tCode = 7
		} else if time == 2 {
			tCode = 9
		} else if time == 3 {
			tCode = 10
		} else {
			tCode = 0
		}
	} else {
		if time == 0 {
			tCode = 28
		} else if time == 1 {
			tCode = 1
		} else if time == 2 {
			tCode = 4
		} else {
			tCode = 0
		}
	}
	if tCode == 0 {
		fmt.Println("Could not get time code.")
	}
	return tCode
}

func getSeatCode(seat string) int {
	url := "http://hi.hana.hs.kr/SYSTEM_Plan/Lib_System/Lib_System_Reservation/popSeat_Reservation.asp"

	rawPayload := fmt.Sprintf("code=001&t_code=&dis_num=%s", seat)
	payload := strings.NewReader(rawPayload)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, _ := client.Do(req)

	defer res.Body.Close()
	// body, _ := ioutil.ReadAll(res.Body)

	// fmt.Println(string(body))
	n, _ := html.Parse(res.Body)
	sCode, ok := getElement(n)
	if !ok {
		fmt.Println("Could not get seat code.")
	}
	sCodeInt, _ := strconv.Atoi(sCode)
	return sCodeInt
}

func getElement(n *html.Node) (string, bool) {
	// spew.Dump(n)
	if n.Type == html.ElementNode && n.Data == "input" {
		// spew.Dump(n)
		for _, a := range n.Attr {
			// spew.Dump(a)
			if a.Key == "name" && a.Val == "s_code" {
				sCode := n.Attr[2].Val
				if sCode != "" {
					return sCode, true
				} else {
					return "", false
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		val, boo := getElement(c)
		if boo {
			return val, true
		}
	}
	return "", false
}
