package matchlib

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type PageParser struct {
	OrignalRequest *http.Request
	OrignalRequestData *[]byte
	Parsers []*PatternParser
}

type PatternParser struct {
	//head|firstLine|(rawPath|rawQuery):[start,end]
	//head|headers|(key|value):[start,end]
	//To Identify
	Where string
	Start int
	End int
	//{Key|Proto:variable|Encoder1|Encoder2|..|EncoderN|Flag}
	Pattern string
	OriginalPattern string
	//The limit qualifier doesn't have to be '{' or '}'.It's identify by Key
	Key string
	Proto string
	Variable string
	Encoders []string
	Flag string
	orignalRequestData *[]byte
}

func (patternParser *PatternParser) parse() error {
	if len(patternParser.Pattern) == 0 {
		patternParser.Pattern = string((*patternParser.orignalRequestData)[patternParser.Start+1:patternParser.End])
	}
	pattern,err := url.QueryUnescape(string(patternParser.Pattern[0:len(patternParser.Pattern)]))
	for ;strings.Contains(pattern,"%"); {
		pattern,err = url.QueryUnescape(pattern)
		if err != nil {
			return nil
		}
	}
	if err != nil {
		return err
	}
	if len(pattern) == 0 {
		return nil
	}

	patternSplit := strings.Split(pattern,"|")
	var key string
	key = patternSplit[0]
	var variable string
	if len(patternSplit) >= 2 {
		variable = patternSplit[1]
	}
	var flag string
	if len(patternSplit) >= 3 {
		flag = patternSplit[len(patternSplit)-1]
	}
	encoders := []string{}
	for i:=2;i < len(patternSplit)-1;i++ {
		encoders = append(encoders,patternSplit[i])
	}
	patternParser.Key = key
	patternParser.Flag = flag

	tmp := strings.Split(variable,":")
	if len(tmp) == 1 {
		patternParser.Proto = "dict"
	} else if len(tmp) == 2 {
		patternParser.Proto = tmp[0]
	}

	patternParser.Variable = tmp[len(tmp)-1]
	patternParser.Encoders = encoders
	return nil
}

func bytesMatchAll(s []byte,substring []byte) []int {
	tmp := s
	lastIndex := 0
	res := []int{}
	i := 0
	for ;bytes.Index(tmp,substring)!=-1; {
		i = bytes.Index(tmp,substring)
		res = append(res,lastIndex+i)
		tmp = tmp[i+len(substring):]
		lastIndex+=len(substring)+i
	}
	return res
}

func isKeyPlaceOk(rawData *[]byte,occurrence int) bool {
	return true
}

type ByFlag [] *PatternParser

func (a ByFlag) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFlag) Less(i, j int) bool { return a[i].Flag < a[j].Flag }
func (a ByFlag) Len() int           { return len(a) }
func MatchRequestRawData(rawData *[]byte) (*PageParser,error) {
	key := []byte("thisismylife")
	keyOccurrence := bytesMatchAll(*rawData,key)
	pageParser := &PageParser{}
	pageParser.OrignalRequestData = rawData
	for _,index := range keyOccurrence {
		if !isKeyPlaceOk(rawData,index) {
			continue
		}
		patternParser := PatternParser{}
		if index == 0 {
			return nil,errors.New("wrong format of request data")
		}
		qualifier := string((*rawData)[index-1:index])
		var unqualifier string
		var isUrlEncodedFlag bool
		if index < 3 || (*rawData)[index-3] != '%' {
			isUrlEncodedFlag = false
			if qualifier == "{" {
				unqualifier = "}"
			} else if qualifier == "(" {
				unqualifier = ")"
			} else if qualifier == "[" {
				unqualifier = "]"
			} else if qualifier == "<" {
				unqualifier = ">"
			} else {
				unqualifier = qualifier
			}
		} else {
			isUrlEncodedFlag = true
			qualifier = string((*rawData)[index-3:index])
			unescapeQualifier,err := url.QueryUnescape(qualifier)
			if err != nil {

			}
			if unescapeQualifier == "{" {
				unqualifier = url.QueryEscape("}")
			} else if unescapeQualifier == "(" {
				unqualifier = url.QueryEscape(")")
			} else if unescapeQualifier == "[" {
				unqualifier = url.QueryEscape("]")
			} else if unescapeQualifier == "<" {
				unqualifier = url.QueryEscape(">")
			} else {
				unqualifier = qualifier
			}
		}
		start := index-1
		end := start + len(key)

		for i:=end;i < len(*rawData)+len(qualifier);i++ {
			if strings.ToLower(string((*rawData)[i:i+len(qualifier)])) == strings.ToLower(unqualifier) {
				end = i
				break
			}
		}

		if isUrlEncodedFlag == true {
			patternParser.OriginalPattern = string((*rawData)[start-2:end+3])
		} else {
			patternParser.OriginalPattern = string((*rawData)[start:end+1])
		}

		patternParser.Start = start
		patternParser.End = end
		patternParser.orignalRequestData = rawData
		patternParser.parse()
		pageParser.Parsers = append(pageParser.Parsers,&patternParser)

		sort.Sort(ByFlag(pageParser.Parsers))
	}
	return pageParser,nil
}

func isValidParamPlace(start int,end int,requestData *[]byte) bool {
	firstLineEnd := bytes.Index(*requestData,[]byte("\r\n"))
	request_uri := bytes.Split((*requestData)[:firstLineEnd],[]byte(" "))[1]
	path := bytes.Split(request_uri,[]byte("?"))[0]
	split := bytes.Split(path,[]byte("."))
	ext := string(split[len(split)-1])
	if ext == ".gif" {
		return false
	}
	if end <= firstLineEnd {
		return true
	}

	bodyStart := bytes.Index(*requestData,[]byte("\r\n\r\n"))
	if start > bodyStart {
		return true
	}

	return false
}

func MatchRequestRawDataPassive(rawData *[]byte) (*PageParser,error) {
	key := []byte("thisismylife")
	keyOccurrence := bytesMatchAll(*rawData,key)
	pageParser := &PageParser{}
	pageParser.OrignalRequestData = rawData
	for _,index := range keyOccurrence {
		if !isKeyPlaceOk(rawData,index) {
			continue
		}
		patternParser := PatternParser{}
		if index == 0 {
			return nil,errors.New("wrong format of request data")
		}
		qualifier := string((*rawData)[index-1:index])
		var unqualifier string
		var isUrlEncodedFlag bool
		if index < 3 || (*rawData)[index-3] != '%' {
			isUrlEncodedFlag = false
			if qualifier == "{" {
				unqualifier = "}"
			} else if qualifier == "(" {
				unqualifier = ")"
			} else if qualifier == "[" {
				unqualifier = "]"
			} else if qualifier == "<" {
				unqualifier = ">"
			} else {
				unqualifier = qualifier
			}
		} else {
			isUrlEncodedFlag = true
			qualifier = string((*rawData)[index-3:index])
			unescapeQualifier,err := url.QueryUnescape(qualifier)
			if err != nil {

			}
			if unescapeQualifier == "{" {
				unqualifier = url.QueryEscape("}")
			} else if unescapeQualifier == "(" {
				unqualifier = url.QueryEscape(")")
			} else if unescapeQualifier == "[" {
				unqualifier = url.QueryEscape("]")
			} else if unescapeQualifier == "<" {
				unqualifier = url.QueryEscape(">")
			} else {
				unqualifier = qualifier
			}
		}
		start := index-1
		end := start + len(key)

		for i:=end;i < len(*rawData)+len(qualifier);i++ {
			if strings.ToLower(string((*rawData)[i:i+len(qualifier)])) == strings.ToLower(unqualifier) {
				end = i
				break
			}
		}

		if !isValidParamPlace(start,end,rawData) {
			return pageParser,nil
		}

		if isUrlEncodedFlag == true {
			patternParser.OriginalPattern = string((*rawData)[start-2:end+3])
		} else {
			patternParser.OriginalPattern = string((*rawData)[start:end+1])
		}

		patternParser.Start = start
		patternParser.End = end
		patternParser.orignalRequestData = rawData
		patternParser.parse()
		pageParser.Parsers = append(pageParser.Parsers,&patternParser)

		sort.Sort(ByFlag(pageParser.Parsers))
	}
	return pageParser,nil
}