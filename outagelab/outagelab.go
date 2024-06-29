package outagelab

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/outagelab/go-sdk/internal/models"
)

type outageLabClient struct {
	options     Options
	accountData *models.Account
}

func NewClient(options Options) *outageLabClient {
	client := &outageLabClient{
		options: options,
	}

	client.init()

	return client
}

type Options struct {
	Application string
	Environment string
	ApiKey      string
	Host        string
}

func (olc *outageLabClient) init() {
	http.DefaultTransport = olc.NewTransport(http.DefaultTransport)
	go olc.pollLoop()
}

func (olc *outageLabClient) pollLoop() {
	for true {
		olc.accountData = getAccountData(&olc.options)
		time.Sleep(5 * time.Second)
	}
}

func getAccountData(options *Options) *models.Account {
	reqJson, err := json.Marshal(map[string]string{
		"application": options.Application,
		"environment": options.Environment,
	})
	if err != nil {
		return nil
	}

	req, err := http.NewRequest("POST", options.Host+"/api/datapage", bytes.NewBuffer(reqJson))
	if err != nil {
		return nil
	}

	req.Header["x-api-key"] = []string{options.ApiKey}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}

	resJson, err := io.ReadAll(res.Body)
	if err != nil {
		return nil
	}

	var account models.Account
	err = json.Unmarshal(resJson, &account)
	if err != nil {
		return nil
	}

	return &account
}

func (olc *outageLabClient) NewTransport(transport http.RoundTripper) http.RoundTripper {
	return newTransport(transport, olc)
}
