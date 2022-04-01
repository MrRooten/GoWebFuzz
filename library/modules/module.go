package modules

import "net/http"

type ModuleMethods interface {
	MatchRequest(req *http.Request)
	MatchResponse(req *http.Request,res *http.Response)
}

type ModuleInfo struct {
	Name string
	Author string
	Description string
	References [][]string

}