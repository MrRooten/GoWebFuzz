package matchproxy

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/google/martian"
	"github.com/google/martian/martianhttp"
	"github.com/google/martian/martianlog"
	"github.com/google/martian/mitm"
	"github.com/google/martian/verify"
	"gowebfuzz/library/fuzz"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

type MatchConfig struct {
	//How many pages will be match
	MatchPages int
	//Match Pattern
	//like ${proto:variables MatchKey flag}
	MatchPattern string
	//Match Key
	MatchKey string

}

type MatchProxy struct {
	//The proxy address
	Addr string
	//The config of the match proxy
	Config MatchConfig
	//The orignal proxy
	Proxy martian.Proxy
	//TLS config
	GenerateCA bool

	Cert string

	Key string
	//downstream-proxy
	DsProxy string

	tr http.Transport
}

type MatchRequestHandler struct {
	MatchPattern string
	Config MatchConfig
}

func (mrh MatchRequestHandler) ModifyRequest(req *http.Request) error {
	if req.Method == "CONNECT" {
		return nil
	}
	fuzz.Drive(req)
	return nil
}


type MatchResponseHandler struct {

}

func (mrh MatchResponseHandler)  ModifyResponse(res *http.Response) error {
	return nil
}


func NewMatchProxy(matchProxy MatchProxy) (*MatchProxy,error) {
	martian.Init()

	p := martian.NewProxy()

	l, err := net.Listen("tcp", matchProxy.Addr)
	if err != nil {
		log.Fatal(err)
	}


	tr := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}
	p.SetRoundTripper(tr)

	if matchProxy.DsProxy != "" {
		u, err := url.Parse(matchProxy.DsProxy)
		if err != nil {
			log.Fatal(err)
		}
		p.SetDownstreamProxy(u)
	}


	var x509c *x509.Certificate
	var priv interface{}

	if matchProxy.Cert != "" && matchProxy.Key != "" {
		tlsc, err := tls.LoadX509KeyPair(matchProxy.Cert, matchProxy.Key)
		if err != nil {
			return nil,err
		}
		priv = tlsc.PrivateKey

		x509c, err = x509.ParseCertificate(tlsc.Certificate[0])
		if err != nil {
			return nil,err
		}
	}

	if x509c != nil && priv != nil {
		mc, err := mitm.NewConfig(x509c, priv)
		if err != nil {
			return nil,err
		}

		mc.SetValidity(time.Hour)
		mc.SetOrganization("Martian")
		mc.SkipTLSVerify(false)

		p.SetMITM(mc)

		tl, err := net.Listen("tcp", ":4443")
		if err != nil {
			log.Fatal(err)
		}

		go p.Serve(tls.NewListener(tl, mc.TLS()))
	}

	p.SetRequestModifier(MatchRequestHandler{MatchPattern: "GET"})
	p.SetResponseModifier(MatchResponseHandler{})

	m := martianhttp.NewModifier()


	logger := martianlog.NewLogger()
	logger.SetDecode(true)




	// Verify assertions.
	vh := verify.NewHandler()
	vh.SetRequestVerifier(m)
	vh.SetResponseVerifier(m)


	// Reset verifications.
	rh := verify.NewResetHandler()
	rh.SetRequestVerifier(m)
	rh.SetResponseVerifier(m)



	go p.Serve(l)
	matchProxy.Proxy = *p
	return &matchProxy,nil
}

func (matchProxy *MatchProxy) CloseProxy() {
	matchProxy.Proxy.Close()
}




