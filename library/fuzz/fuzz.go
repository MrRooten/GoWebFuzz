package fuzz

import (
	"gowebfuzz/library/fuzz/matchlib"
	"gowebfuzz/library/network/gwfhttp"
	"gowebfuzz/library/utils"
	"net/http"
)
func StartFuzzing(pageParser *matchlib.PageParser) {

}
func Drive(requeset *http.Request) error {
	reqp := gwfhttp.RequestPacket{}
	reqp.ReadFromHTTPRequest(requeset)
	pageParser,err := matchlib.MatchRequest(requeset)
	if err != nil {
		utils.Log.Error(err.Error())
	}
	return nil
}