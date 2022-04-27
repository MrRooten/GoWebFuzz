package information

import (
	"gowebfuzz/modules"
	"net/http"
)

type FingerPrint struct {
	modules.ModuleInfo
}

type IdentityPrint struct {
	Type    string
	Feature string
}

func (fp FingerPrint) MatchRequest(req *http.Request) error {
	return nil
}

func (fp FingerPrint) MatchResponse(req *http.Request, res *http.Response, body *[]byte) error {
	return nil
}
