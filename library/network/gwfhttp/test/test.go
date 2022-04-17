package main

import (
	"crypto/tls"
	"fmt"
	"gowebfuzz/library/network/gwfhttp"
	"io/ioutil"
	"net/http"
	"time"

	"os"
	"strings"
)

func httpsTest() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	request,err := http.NewRequest("GET","https://www.bing.com",strings.NewReader(""))
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	req := &http.Client{
		Transport: tr,
	}
	a,err := req.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	bs ,err := ioutil.ReadAll(a.Body)
	fmt.Println(string(bs))
}

func burpfileRequest() {
	reqp := gwfhttp.RequestPacket{}
	//reqp.ReadFromHTTPRequest(request)
	file,_ := os.Open(".\\burpsuite")
	bs,_ := ioutil.ReadAll(file)
	reqp.ReadFromBytes(&bs,true)
	reqp.SetHttpConfig(time.Duration(5*time.Second), "", 5)
	for i:=0;i < 10;i++ {
		//reqp.ReadFromBytes(&bs, true)

		a,err := reqp.SendPacket()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(a.GetBody(-1))
		fmt.Println(i)
	}
}
func main() {
	str, _ := os.Getwd()
	fmt.Println(str)
	burpfileRequest()
	//client := &http.Client{
	//	Transport: &http.Transport{
	//		MaxIdleConnsPerHost: 20,
	//	},
	//	Timeout: 10 * time.Second,
	//}
	//for i:=0;i < 10;i++ {
	//	request, _ := http.NewRequest("POST", "http://baidu.com", strings.NewReader("abcefg"))
	//	request.Header.Add("Accept-Encoding", "identity")
	//	res,_ := client.Do(request)
	//	io.Copy(ioutil.Discard, res.Body)
	//	res.Body.Close()
	//	fmt.Println(i)
	//}

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
