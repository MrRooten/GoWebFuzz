package information

import (
	"gowebfuzz/library/modules"
	"net/http"
)

type FingerPrint struct {
	modules.ModuleInfo
}

func (fp FingerPrint)MatchRequeset(req *http.Request) error {
	return nil
}

func (fp FingerPrint)MatchResponse(req *http.Request,res *http.Response,body *[]byte) error {
	return nil
}