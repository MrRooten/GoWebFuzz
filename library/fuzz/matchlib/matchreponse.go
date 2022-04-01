package matchlib

import (
	"net/http"
)

func MatchResponseHTTPResponse(r *http.Response) bool {
	return true
}

type ResponseMatch struct {

}

type Response struct {

}

func MatchResponse(req *http.Request,resp *http.Response) {

}
