package information

import (

	"fmt"
	"gowebfuzz/library/modules"

	"net/http"
)

type Test struct {
	modules.ModuleInfo
}

func (test Test)MatchRequest(req *http.Request) error {
	return nil
}

func (test Test)MatchResponse(req *http.Request, res *http.Response,body *[]byte) error {
	fmt.Println(req)
	fmt.Println(res)
	fmt.Println(string(*body))
	return nil
}