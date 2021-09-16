package busterlib

import (
	"bufio"
	"bytes"
	"fmt"
	"gowebfuzz/library/fuzz/matchlib"
	"io"
	"os"
)

func replaceBytes(bytesToBeReplaced *[]byte,bytesToReplace []byte,start int,end int) *[]byte {
	tmp := (*bytesToBeReplaced)[end:]
	*bytesToBeReplaced = append((*bytesToBeReplaced)[:start],bytesToReplace...)
	*bytesToBeReplaced = append(*bytesToBeReplaced,tmp...)
	return bytesToBeReplaced
}

func parseProtoVariableToDict(dicts *[][]string,parser *matchlib.PatternParser) {
	dict := []string{}
	proto := parser.Proto
	if proto == "file" {
		fi, err := os.Open(parser.Variable)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		defer fi.Close()
		br := bufio.NewReader(fi)
		for {
			a, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			dict = append(dict,string(a))
		}
	} else if proto == "dict" {
		if parser.Variable == "//usernames" {
			dict = []string{
				"username1",
				"username2",
			}
		} else if parser.Variable == "//passwords" {
			dict = []string {
				"password1",
				"password2",
				"password3",
			}
		}
	}
	*dicts = append(*dicts,dict)
}

func isIteratorDone(mp *map[int][]int) bool {
	if (*mp)[0][0] == (*mp)[0][1] {
		return true
	}
	return false
}

func replace(args *[]string,bytesToSend *[]byte,page *matchlib.PageParser) *[]byte {
	for index,patternParser := range page.Parsers {
		*bytesToSend = bytes.Replace(*bytesToSend,[]byte(patternParser.OriginalPattern),[]byte((*args)[index]),-1)
	}
	return nil
}

func iterator(dicts* [][]string,page *matchlib.PageParser) {
	mi := map[int][]int {}
	for index,dict := range *dicts {
		mi[index] = []int{len(dict),0}
	}

	for {
		sendRequestData := *page.OrignalRequestData
		if isIteratorDone(&mi) {
			break
		}
		//Process place
		//fmt.Println((*dicts)[0][mi[0][1]],":",(*dicts)[1][mi[1][1]])
		simplize := func (index int) string {
			return (*dicts)[index][(mi)[index][1]]
		}
		//fmt.Println(simplize(0),":",simplize(1))
		argsList := []string{}
		for index := range *dicts {
			argsList = append(argsList,simplize(index))
		}
		replace(&argsList,&sendRequestData,page)
		fmt.Println(string(sendRequestData))
		lastKey:=len(mi)-1
		mi[lastKey][1]++
		index := lastKey
		for ;index>0;index-- {
			if mi[index][1] == mi[index][0] {
				mi[index][1] = 0
				mi[index-1][1]++
			}
		}

	}
}


func ProcIterator(pageParser *matchlib.PageParser) {
	dicts := [][]string {}
	for _,parser := range pageParser.Parsers {
		parseProtoVariableToDict(&dicts,parser)
	}

	iterator(&dicts,pageParser)
}
