package outagelab

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/outagelab/go-sdk/internal/models"
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
	rule := olt.getOutageRule(req)

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

func (olt *outageLabTransport) getOutageRule(req *http.Request) *models.OutageRule {
	client := olt.outageLabClient
	options := client.options

	account := client.accountData
	if account == nil {
		return nil
	}

	var application *models.Application
	for _, app := range account.Applications {
		if app.ID == options.Application {
			application = app
		}
	}

	if application == nil {
		return nil
	}

	var environment *models.Environment
	for _, env := range application.Environments {
		if env.ID == options.Environment && env.Enabled {
			environment = env
		}
	}

	if environment == nil {
		return nil
	}

	var outageRule *models.OutageRule
	for _, rule := range application.Rules {
		if rule.Host == req.Host && rule.Enabled {
			outageRule = rule
		}
	}

	if outageRule == nil {
		return nil
	}

	return outageRule
}
