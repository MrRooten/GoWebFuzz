package main

import (
	"fmt"
	"gowebfuzz/library/network/gwfhttp"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	str, _ := os.Getwd()
	fmt.Println(str)
	request, _ := http.NewRequest("POST", "http://baidu.com", strings.NewReader("abcefg"))
	request.Header.Add("Accept-Encoding", "identity")
	request.Header.Add("Host","http://google.com")
	reqp := gwfhttp.RequestPacket{}
	//reqp.ReadFromHTTPRequest(request)
	file,_ := os.Open(".\\burpsuite")
	bs,_ := ioutil.ReadAll(file)
	reqp.ReadFromBytes(&bs,false)
	fmt.Println(reqp)
	reqp.SetHttpConfig(time.Duration(5*time.Second),"http://127.0.0.1:8080",1)
	fmt.Println(string(reqp.WriteToBytes()))
	resp,err := reqp.SendPacket()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(resp.GetBody(1024))
	//request, err := http.NewRequest("POST", "", strings.NewReader("abcefg"))
	//request.URL.Host = "baidu.com"
	//request.URL.Path = "/%uabc"
	//request.URL.Scheme = "http"
	//client := &http.Client{}
	//res,_ := client.Do(request)
	//fmt.Printf("%#v\n",err)
	//a,_ := ioutil.ReadAll(res.Body)
	//fmt.Printf("%#v\n", string(a))
}
