package main

import (
	"fmt"
	"gowebfuzz/library/network/matchproxy"
)

func main() {
	mp := matchproxy.MatchProxy{Addr: ":9999",
		Cert: "D:\\wsl\\installed\\kali-linux\\rootfs\\home\\sam0ple\\.mitmproxy\\mitmproxy-ca-cert.cer",
		Key: "D:\\wsl\\installed\\kali-linux\\rootfs\\home\\sam0ple\\.mitmproxy\\mitmproxy-ca.pem",
		DsProxy: "http://127.0.0.1:8080"}

	mp1,err := matchproxy.NewMatchProxy(mp)

	if err != nil {
		fmt.Println(err)
	}
	var a string
	fmt.Scanf("%s",a)
	defer mp1.CloseProxy()
}


