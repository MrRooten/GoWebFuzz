package information

import (
	"gowebfuzz/library/modules"
	"gowebfuzz/library/network/gwfhttp"
	"net/http"
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

var Records []*Pair
var lock sync.Mutex

type HttpHistory struct {
	modules.ModuleInfo
}

func (history HttpHistory) MatchRequest(req *http.Request) error {
	lock.Lock()
	defer lock.Lock()

	return nil
}

func (history HttpHistory) MatchResponse(req *http.Request, res *http.Response, body *[]byte) error {
	return nil
}
