package matchlib

import (
	"bytes"
	"errors"
	"net/http"
)

type PageParser struct {
	orignalRequest *http.Request
	orignalRequestData []byte
	patternParsers []*PatternParser
}

type PatternParser struct {
	//head|firstLine|(rawPath|rawQuery):[start,end]
	//head|headers|(key|value):[start,end]
	//To Identify
	Where string
	Start int
	End int
	//{Key|Proto:variable|Encoder1|Encoder2|..|EncoderN|Flag}
	//The limit qualifier doesn't have to be '{' or '}'.It's identify by Key
	Key string
	Proto string
	variable string
	Encoders []string
	Flag string
}

func MatchRequestHTTPRequest(r *http.Request) (*PageParser,error) {

}

func bytesMatchAll(s []byte,substring []byte) []int {

}

func MatchRequestRawData(rawData *[]byte) (*PageParser,error) {
	key := []byte("abcdefg")
	keyOccurrence := bytesMatchAll(*rawData,key)
	pageParser := &PageParser{}
	for _,index := range keyOccurrence {
		patternParser := PatternParser{}
		if index == 0 {
			return nil,errors.New("wrong format of request data")
		}
		qualifier := (*rawData)[index-1]
		var unqualifier byte
		if qualifier == '{' {
			unqualifier = '}'
		} else if qualifier == '(' {
			unqualifier = ')'
		} else if qualifier == '[' {
			unqualifier = ']'
		} else if qualifier == '<' {
			unqualifier = '>'
		} else {
			unqualifier = qualifier
		}

		start := index - 1
		end := start + len(key)
		for i:=start;i < len(*rawData);i++ {
			if (*rawData)[i] == unqualifier {
				end = i
			}
		}

		patternParser.Start = start
		patternParser.End = end

		
		pageParser.patternParsers = append(pageParser.patternParsers,&patternParser)
	}
}
