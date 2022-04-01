package main

import (
	"fmt"
	"gowebfuzz/library/fuzz"
	"gowebfuzz/library/network/gwfhttp"
	"net/http"
	"time"
)

func testCase() {
	fz := fuzz.Fuzz{}
	fz.Init(10,"go://0,go://1,go://2,go://3","",6)
	var testCase [4]*http.Request
	testCase[0], _ = http.NewRequest("GET","http://baidu1.com",nil)
	testCase[1], _ = http.NewRequest("GET","http://baidu2.com",nil)
	testCase[2], _ = http.NewRequest("GET","http://baidu1.com",nil)
	testCase[3], _ = http.NewRequest("GET","http://baidu2.com",nil)

	go fz.StartFuzzing()
	for _,request := range testCase {

		fmt.Println(request)
		fz.Drive(request)
	}

	time.Sleep(10*time.Second)
}

func testQueue() {
	queue := fuzz.ReuqestQueue{}
	var testCase [2]*http.Request
	testCase[0], _ = http.NewRequest("GET","http://baidu.com",nil)
	testCase[1], _ = http.NewRequest("GET","http://baidu.com",nil)
	for _,request := range testCase {
		b := gwfhttp.RequestPacket{}
		b.ReadFromHTTPRequest(request)
		queue.Push(&b)
	}

	a := queue.Take()
	a = queue.Take()
	a = queue.Take()
	fmt.Println(a)
}

func main() {
	testCase()
}