package modules

import "net/http"

type ModuleMethods interface {
	MatchRequest(req *http.Request) error
	MatchResponse(req *http.Request,res *http.Response) error
}

type ModuleInfo struct {
	Name string
	Author string
	Description string
	References [][]string

}