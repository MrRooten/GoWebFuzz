package main

import (
	"fmt"
	"gowebfuzz/library/fuzz/busterlib"
	"gowebfuzz/library/fuzz/matchlib"
	"io/ioutil"
	"os"
)

func isIteratorDone(mp *map[int][]int) bool {
	if (*mp)[0][0] == (*mp)[0][1] {
		return true
	}
	return false
}
func iterator() {
	mi := map[int][]int {}

	usernames := []string{
		"username1",
		"username2",
		"username3",
	}

	passwords := []string{
		"password1",
		"password2",
		"password3",
	}

	passwords2 := []string{
		"password21",
		"password22",
		"password23",
	}

	mi[0] = []int{len(usernames),0}
	mi[1] = []int{len(passwords),0}
	mi[2] = []int{len(passwords2),0}
	for ;; {
		if isIteratorDone(&mi) {
			break
		}
		fmt.Println(usernames[mi[0][1]],":",passwords[mi[1][1]],":",passwords2[mi[2][1]])

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

func replaceBytes(bytesToBeReplaced *[]byte,bytesToReplace []byte,start int,end int) *[]byte {
	tmp := (*bytesToBeReplaced)[end:]
	*bytesToBeReplaced = append((*bytesToBeReplaced)[:start],bytesToReplace...)
	*bytesToBeReplaced = append(*bytesToBeReplaced,tmp...)
	return bytesToBeReplaced
}
func main() {
	//s := []byte("%3babcdefg|hij|base64|3123%3b")
	//a,_ := matchlib.MatchRequestRawData(&s)
	//fmt.Printf("%#v\n",a.Parsers[0])
	//target := []byte("abcdefg")
	//replaceBytes(&target,[]byte("123"),0,3)
	//fmt.Println(string(target))

	file,_ := os.Open(".\\burpsuite")
	bs,_ := ioutil.ReadAll(file)
	p,_ := matchlib.MatchRequestRawData(&bs)
	//fmt.Printf("%#v",p.Parsers)
	busterlib.ProcIterator(p)
}
