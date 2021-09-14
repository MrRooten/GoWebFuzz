package matchlib

import "net/http"

func MatchResponseHTTPResponse(r *http.Response) bool {
	return true
}
