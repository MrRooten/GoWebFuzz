package matchlib

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Modifier struct {
	Body *[]byte
	Req  *http.Request
	Pair *Pair
}

func NewModifier() *Modifier {
	modifier := &Modifier{}
	modifier.Req = new(http.Request)
	modifier.Pair = new(Pair)

	return modifier
}

type LogLevel int

const (
	DbgLevel   LogLevel = 0
	InfoLevel  LogLevel = 1
	WarnLevel  LogLevel = 2
	ErrorLevel LogLevel = 3
	KeyLevel   LogLevel = 4
)

type Message struct {
	Provider   string
	TimeCreate string
	Message    string
	Level      LogLevel
}

type MessageListener struct {
	MessageChannel chan *Message
	LevelToLog     LogLevel
}

func (listener *MessageListener) ListenEvent() {
	for {
		message := <-listener.MessageChannel
		if message.Level >= listener.LevelToLog {
			fmt.Println(message.Message)
		}
	}
}

func (modifier Modifier) ModifyRequest(req *http.Request) error {
	*modifier.Req = *req
	Records.AddReq(modifier.Pair, req)
	return nil
}

func (modifier Modifier) ModifyResponse(res *http.Response) error {
	if res.ContentLength > 0x10000 {
		return nil
	}
	_tmp, _ := ioutil.ReadAll(res.Body)
	modifier.Body = &_tmp
	res.Body = ioutil.NopCloser(bytes.NewReader(*modifier.Body))
	var err error
	go func() error {
		if modifier.Req.Method == "CONNECT" {
			return nil
		}

		if strings.Contains(res.Header.Get("Content-Type"), "gzip") {
			reader, err := gzip.NewReader(bytes.NewReader(*modifier.Body))
			if err != nil {
				return nil
			}
			_tmp, err = ioutil.ReadAll(reader)
			modifier.Body = &_tmp
			if err != nil {
				return nil
			}
		} else if strings.Contains(modifier.Req.Header.Get("Accept-Encoding"), "gzip") {
			reader, err := gzip.NewReader(bytes.NewReader(*modifier.Body))
			if err != nil {
				return nil
			}
			_tmp, err = ioutil.ReadAll(reader)
			modifier.Body = &_tmp
			if err != nil {
				return nil
			}
		}

		Records.AddRes(modifier.Pair, res)
		modifier.Pair.Response.Body = *modifier.Body
		//finder := information.LinksFinder{}
		//err := finder.MatchResponse(mrh.Req, res, &mrh.Body)

		if err != nil {
			return nil
		}
		return nil
	}()
	fmt.Println(len(Records.Pairs))
	return err
}
