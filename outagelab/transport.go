package outagelab

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

type outageLabTransport struct {
	wrappedTransport http.RoundTripper
	outageLabClient  *outageLabClient
}

func newTransport(
	wrappedTransport http.RoundTripper,
	outageLabClient *outageLabClient,
) *outageLabTransport {
	return &outageLabTransport{
		wrappedTransport: wrappedTransport,
		outageLabClient:  outageLabClient,
	}
}

func (olt *outageLabTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rule := olt.getHttpClientOutageRule(req)

	if rule == nil {
		return olt.wrappedTransport.RoundTrip(req)
	}

	if rule.Duration > 0 {
		time.Sleep(time.Duration(rule.Duration) * time.Second)
	}

	if rule.Status == 0 {
		return olt.wrappedTransport.RoundTrip(req)
	}

	res := &http.Response{
		Status:        fmt.Sprintf("%v %v", rule.Status, http.StatusText(rule.Status)),
		StatusCode:    rule.Status,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          io.NopCloser(bytes.NewBuffer([]byte{})),
		ContentLength: 0,
		Request:       req,
		Header:        make(http.Header, 0),
	}

	return res, nil
}

func (olt *outageLabTransport) getHttpClientOutageRule(req *http.Request) *httpClientRequestV1 {
	client := olt.outageLabClient

	datapage := client.datapage
	if datapage == nil {
		return nil
	}

	var outageRule *httpClientRequestV1
	for _, rule := range datapage.Rules {
		switch rule.Type {
		case "http-client-request.v1":
			rule := rule.HttpClientRequestV1
			if rule.Host == req.Host {
				outageRule = rule
			}
		}
	}

	if outageRule == nil {
		return nil
	}

	return outageRule
}
