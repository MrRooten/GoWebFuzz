package main

import (
	"fmt"
	"gowebfuzz/library/network/matchproxy"
)

func main() {
	mp := matchproxy.MatchProxy{Addr: ":9999",
		Cert:    "./mitmproxy-ca-cert.cer",
		Key:     "./mitmproxy-ca.pem",
		DsProxy: ""}

	mp1, err := matchproxy.NewMatchProxy(mp)

	if err != nil {
		fmt.Println(err)
	}
	var a string
	fmt.Scanf("%s", a)
	defer mp1.CloseProxy()
}
