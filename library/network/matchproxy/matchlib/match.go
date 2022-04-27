package matchlib

import (
	"bytes"
	"compress/gzip"
	"errors"
	"gowebfuzz/library/fuzz"
	"gowebfuzz/library/network/gwfhttp"
	information2 "gowebfuzz/modules/information"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

type HttpHistoryConfig struct {
	MaxItems    int
	MaxItemSize int
}

type Pair struct {
	RecordId int
	Request  gwfhttp.RequestPacket
	Response gwfhttp.ResponsePacket
}

var config HttpHistoryConfig

type _Records struct {
	Pairs        []*Pair
	LatestId     int
	Size         int
	MaxItemsSize int
	locker       sync.Mutex
}

func (records *_Records) AddReq(pair *Pair, req *http.Request) {
	records.locker.Lock()
	defer records.locker.Unlock()
	err := pair.Request.ReadFromHTTPRequest(req)
	if err != nil {
		return
	}
	if len(records.Pairs) >= 1000 {
		records.Pairs = records.Pairs[len(records.Pairs)-1000 : len(records.Pairs)]
	}
	records.LatestId = records.LatestId + 1
	pair.RecordId = records.LatestId
	//records.Pairs[pair.RecordId] = pair
	records.Pairs = append(records.Pairs, pair)
}

func (records *_Records) AddRes(pair *Pair, res *http.Response) {
	records.locker.Lock()
	defer records.locker.Unlock()
	err := pair.Response.ReadFromHTTPResponse(res)
	if err != nil {
		return
	}

}

var Records _Records

// linkfinder information.LinksFinder

var sensitiveLeak information2.SensitiveiLeak
var siteMap information2.SiteMap

func Exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

type RequestTuple struct {
	Request   *http.Request
	Response  *http.Response          //Reverve for later update
	ResPacket *gwfhttp.ResponsePacket //Reverve for later update
}

type MatchRequestHandler struct {
	MatchPattern string
	Fuzz         *fuzz.Fuzz
	Request      *http.Request
	Pair         *Pair
}

var locker sync.Mutex

func (mrh MatchRequestHandler) ModifyRequest(req *http.Request) error {
	if req.Method == "CONNECT" {
		return nil
	}
	*mrh.Request = *req
	*mrh.Pair = Pair{}
	Records.AddReq(mrh.Pair, req)
	//Save the requeset and response to a queue
	//mrh.Fuzz.Drive(req)

	return nil
}

type MatchResponseHandler struct {
	Pair *Pair
	Req  *http.Request
	Body *[]byte
}

func (mrh MatchResponseHandler) ModifyResponse(res *http.Response) error {
	if res.ContentLength > 0x10000 {
		return nil
	}
	_tmp, _ := ioutil.ReadAll(res.Body)
	mrh.Body = &_tmp
	res.Body = ioutil.NopCloser(bytes.NewReader(*mrh.Body))
	var err error
	go func() error {
		if mrh.Req.Method == "CONNECT" {
			return nil
		}

		if strings.Contains(res.Header.Get("Content-Type"), "gzip") {
			reader, err := gzip.NewReader(bytes.NewReader(*mrh.Body))
			if err != nil {
				return nil
			}
			_tmp, err = ioutil.ReadAll(reader)
			mrh.Body = &_tmp
			if err != nil {
				return nil
			}
		} else if strings.Contains(mrh.Req.Header.Get("Accept-Encoding"), "gzip") {
			reader, err := gzip.NewReader(bytes.NewReader(*mrh.Body))
			if err != nil {
				return nil
			}
			_tmp, err = ioutil.ReadAll(reader)
			mrh.Body = &_tmp
			if err != nil {
				return nil
			}
		}

		Records.AddRes(mrh.Pair, res)
		mrh.Pair.Response.Body = *mrh.Body
		//finder := information.LinksFinder{}
		//err := finder.MatchResponse(mrh.Req, res, &mrh.Body)

		if err != nil {
			return nil
		}
		return nil
	}()

	//fmt.Println(len(Records.Pairs))

	return err
}

func NewMatchRequestHandler(routineSize int) *MatchRequestHandler {
	mrh := MatchRequestHandler{}
	mrh.Fuzz = new(fuzz.Fuzz)
	mrh.Fuzz.Init(5, "go://1", "", 5)
	go mrh.Fuzz.StartFuzzing()
	return &mrh
}

func NewMatchResponseHandler(requestHandler *MatchRequestHandler) *MatchResponseHandler {
	mrh := MatchResponseHandler{}
	mrh.Req = requestHandler.Request
	mrh.Pair = requestHandler.Pair
	return &mrh
}
