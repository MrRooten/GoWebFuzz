package fuzz

import (
	"fmt"
	"gowebfuzz/library/fuzz/matchlib"
	"gowebfuzz/library/network/gwfhttp"
	"gowebfuzz/library/utils"
	"net/http"
)
func StartFuzzing(pageParser *matchlib.PageParser) {

}

//Passive fuzz only match uri body and POST body
func Drive(requeset *http.Request) error {
	reqp := gwfhttp.RequestPacket{}
	reqp.ReadFromHTTPRequest(requeset)
	bs := reqp.WriteToBytes()
	pageParser,err := matchlib.MatchRequestRawDataPassive(&bs)
	if err != nil {
		utils.Log.Error(err.Error())
	}

	if pageParser != nil && len(pageParser.Parsers) != 0 {
		fmt.Println(string(bs))
	}
	return nil
}