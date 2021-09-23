package lualib

import (
	lua "github.com/yuin/gopher-lua"
	"gowebfuzz/library/network/gwfhttp"
	"net/http"
)

func RegisterRequestTable(L *lua.LState,r *http.Request) error {
	reqp := gwfhttp.RequestPacket{}
	reqp.ReadFromHTTPRequest(r)
	mt := L.NewTypeMetatable("requeset")

	scheme := r.URL.Scheme
	L.SetField(mt,"scheme",lua.LString(scheme))

	proto := lua.LString(reqp.GetProto())
	L.SetField(mt,"proto",proto)

	requestURI := lua.LString(reqp.GetURI())
	L.SetField(mt,"request_uri",requestURI)

	headers := L.NewTable()
	rawHeaders := reqp.GetHeaders()
	for k,v := range rawHeaders {
		headers.RawSetString(k,lua.LString(v))
	}
	L.SetField(mt,"headers",headers)

	rawBody,err := reqp.GetBody()
	if err != nil {
		return err
	}
	L.SetField(mt,"body",lua.LString(rawBody))
	L.SetGlobal("req",mt)
	return nil
}

func RegisterResponseTable(L *lua.LState,r *http.Response) error {
	resp := gwfhttp.ResponsePacket{}
	resp.ReadFromHTTPResponse(r)
	mt := L.NewTypeMetatable("response")

	proto := lua.LString(resp.GetProto())
	L.SetField(mt,"proto",proto)

	code := lua.LNumber(float64(resp.GetStatusCode()))
	L.SetField(mt,"status_code",code)

	headers := L.NewTable()
	rawHeaders := resp.GetHeaders()
	for k,v := range rawHeaders {
		headers.RawSetString(k,lua.LString(v))
	}
	L.SetField(mt,"headers",headers)


	rawBody,err := resp.GetBody(1024*1024)
	if err != nil {
		return err
	}
	L.SetField(mt,"body",lua.LString(rawBody))
	L.SetGlobal("res",mt)
	return nil
}

var exports = map[string]lua.LGFunction{

}

