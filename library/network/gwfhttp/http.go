package gwfhttp

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type RequestConfig struct {
	Timeout time.Duration
	Proxy string
	FollowDirectionDepth int
}
type RequestPacket struct {
	// Save the orginal data
	rawData []byte
	// Is ssl request or not
	isSSL bool
	// Save Origninal http request
	orignalRequest http.Request
	// GET,POST,PUT,OPTIONS,CONNECT
	httpMethod string
	// /abc/efg
	uri string
	// HTTP/1.1 HTTP/2.0
	httpProtocol string
	// https://example.com/abc/efg
	requestURL string
	// Map of headers
	headers map[string]string

	host string

	httpScheme string

	body []byte

	client http.Client
}

func (request *RequestPacket) SetHttpConfig(timeout time.Duration,proxyString string,redirectDepth int) error{
	request.client.Timeout = timeout

	tr := &http.Transport{}
	if proxyString != "" {
		proxy,err := url.Parse(proxyString)
		if err != nil {
			return err
		}
		tr = &http.Transport{
			Proxy: http.ProxyURL(proxy),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		request.client.Transport = tr
	}

	request.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) > redirectDepth {
			return errors.New("too many redirections")
		}
		return nil
	}

	return nil
}

func (request *RequestPacket) ReadFromHTTPRequest(req *http.Request) error {
	if req == nil {
		return errors.New("can not pass a nil as an argument")
	}

	if req.TLS == nil {
		request.isSSL = false
	} else {
		request.isSSL = true
	}

	request.httpProtocol = ""
	if request.isSSL == true {
		request.httpScheme = "https://"
	} else {
		request.httpScheme = "http://"
	}
	request.orignalRequest = *req
	request.httpMethod = req.Method
	request.uri = req.RequestURI
	request.httpProtocol = req.Proto
	request.requestURL = request.httpScheme + req.Host + "/" +request.uri
	request.headers = map[string]string{}
	for k,v := range req.Header {
		request.headers[k] = strings.Join(v,"")
	}
	return nil
}

func (request *RequestPacket) WriteToBytes() []byte {
	res := []byte("")
	res = append(res[:],[]byte(request.httpMethod)[:]...)
	res = append(res[:],[]byte(" ")...)
	res = append(res[:],[]byte(request.uri)[:]...)
	res = append(res[:],[]byte(" ")...)
	res = append(res[:],[]byte(request.httpProtocol)[:]...)
	res = append(res[:],[]byte("\r\n")...)
	for k,v := range request.headers {
		s := []byte("")
		s = append(s[:],[]byte(k)...)
		s = append(s[:],[]byte(" : ")...)
		s = append(s[:],[]byte(v)...)
		s = append(s[:],[]byte("\r\n")...)
		res = append(res[:],s...)
	}
	res = append(res[:],[]byte("\r\n")...)
	res = append(res[:],request.body...)
	return res
}

func (request *RequestPacket) ReadFromBytes(rawData *[]byte,isSSL bool) error {
	seperatorIndex := bytes.Index(*rawData,[]byte("\r\n\r\n"))
	head := (*rawData)[0:seperatorIndex]
	body := (*rawData)[seperatorIndex+4:len(*rawData)]
	if request.headers == nil {
		request.headers = map[string]string{}
	}

	headSplit := bytes.Split(head,[]byte("\r\n"))
	firstLineSplit := bytes.Split(headSplit[0],[]byte(" "))
	if len(firstLineSplit) != 3 {
		return errors.New("the request is not valid in first line")
	}

	request.httpMethod = string(bytes.Trim(firstLineSplit[0]," "))
	request.uri = string(bytes.Trim(firstLineSplit[1]," "))
	request.httpProtocol = string(bytes.Trim(firstLineSplit[2]," "))
	for i:=1;i < len(headSplit);i++ {
		colonIndex := bytes.Index(headSplit[i],[]byte(":"))
		if colonIndex == -1 {
			return errors.New("the header is not Valid in " + string(headSplit[i]))
		}
		key := bytes.Trim(headSplit[i][0:colonIndex]," ")
		value := bytes.Trim(headSplit[i][colonIndex+1:len(headSplit[i])]," ")
		request.headers[string(key)] = string(value)
		if strings.ToLower(string(key)) == "host" {
			request.host = string(value)
		}
	}

	request.body = body

	request.isSSL = isSSL
	if isSSL == false {
		request.httpScheme = "http://"
	} else {
		request.httpScheme = "https://"
	}

	request.requestURL = request.httpScheme + request.host + request.uri
	return nil
}


func (request *RequestPacket) GetHTTPMethod() string {
	if request.httpMethod != "" {
		return request.httpMethod
	} else if !reflect.DeepEqual(request.orignalRequest,http.Request{}) {
		request.httpMethod = request.orignalRequest.Method
		return request.httpMethod
	} else if len(request.rawData) != 0 {
		request.httpMethod = string(bytes.Split(request.rawData,[]byte(" "))[0])
		return request.httpMethod
	}
	return ""
}

func (request *RequestPacket) GetURI() string {
	if request.uri != "" {
		return request.uri
	} else if !reflect.DeepEqual(request.orignalRequest,http.Request{}) {
		request.uri = request.orignalRequest.RequestURI
		return request.uri
	} else if len(request.rawData) != 0 {
		request.uri = string(bytes.Split(request.rawData,[]byte(" "))[1])
		return request.uri
	}
	return ""
}

func (request *RequestPacket) GetHTTPProtocol() string {
	return request.httpProtocol
}

func (request *RequestPacket) GetRequestURL() string {
	return request.requestURL
}

func (request *RequestPacket) GetHeaders() map[string]string {
	return request.headers
}

func (request *RequestPacket) GetBody() ([]byte,error) {
	if len(request.body) != 0 {
		return request.body,nil
	} else if !reflect.DeepEqual(request.orignalRequest,http.Request{}) {
		reader := request.orignalRequest.Body
		if reader == nil {
			return []byte(""),nil
		}
		body,err := ioutil.ReadAll(reader)
		if err != nil {
			return []byte(""),err
		}
		return body,err
	}  else if len(request.rawData) != 0 {
		body := request.rawData[bytes.Index(request.rawData,[]byte("\r\n")):len(request.rawData)]
		return body,nil
	}
	return []byte(""),errors.New("not valid request bytes")
}

func (request *RequestPacket) SendPacket() (*ResponsePacket,error){
	body,err := request.GetBody()
	if err != nil {
		return nil,err
	}
	req,err := http.NewRequest(request.httpMethod,request.requestURL,bytes.NewReader(body))
	//req.Proto = request.httpProtocol
	for k,v := range request.headers {
		req.Header.Set(k,v)
	}

	client := request.client

	res,err := client.Do(req)

	if err != nil {
		return nil,err
	}
	response := ResponsePacket{}
	response.ReadFromHTTPResponse(res)
	return &response,nil
}


type ResponsePacket struct {
	rawData []byte
	bodyReader io.Reader
	headers map[string]string
	statusCode int
	httpProtocol string
}

func (response *ResponsePacket) ReadFromHTTPResponse(res *http.Response) (*ResponsePacket,error) {
	response.bodyReader = res.Body
	response.headers = map[string]string{}
	for k,v := range res.Header {
		response.headers[k] = strings.Join(v,"")
	}
	response.statusCode = res.StatusCode
	response.httpProtocol = res.Proto
	return response,nil
}

func (response *ResponsePacket) GetBody(size int) ([]byte,error) {
	var body []byte
	if size >= 0 {
		body = make([]byte, size)
	} else {
		return ioutil.ReadAll(response.bodyReader)
	}

	response.bodyReader.Read(body)
	return body,nil
}

func (response *ResponsePacket) GetHeaders() map[string]string {
	return response.headers
}

func (response *ResponsePacket) GetStatusCode() int {
	return response.statusCode
}

func GetBytesFromHTTPRequest(request *http.Request) ([]byte,error) {
	return []byte(""),nil
}