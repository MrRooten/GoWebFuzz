package information

import (
	"gowebfuzz/modules"
	"net/http"
)

type SensitiveiLeak struct {
	modules.ModuleInfo
}

func matchEmails(body *[]byte) []string {
	return nil
}

func matchNumber(body *[]byte) []string {
	return nil
}

func (sl SensitiveiLeak) MatchRequest(req *http.Request) error {
	return nil
}

func (sl SensitiveiLeak) MatchResponse(req *http.Request, res *http.Response, body *[]byte) error {
	return nil
}
