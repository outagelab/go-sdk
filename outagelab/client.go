package outagelab

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type outageLabClient struct {
	options            Options
	datapage           *datapageResponse
	outagelabTransport *outageLabTransport
	cancelPoll         context.CancelFunc
}

func newClient(options Options) *outageLabClient {
	return &outageLabClient{
		options: options,
	}
}

func (olc *outageLabClient) start() {
	ctx, cancel := context.WithCancel(context.Background())
	olc.cancelPoll = cancel

	olc.outagelabTransport = olc.NewTransport(http.DefaultTransport)
	http.DefaultTransport = olc.outagelabTransport

	go olc.pollLoop(ctx)
}

func (olc *outageLabClient) stop() {
	// this will definitely need to be revisited,
	// not safe if other tools are modifying DefaultTransport
	http.DefaultTransport = olc.outagelabTransport.wrappedTransport
	olc.cancelPoll()
}

func (olc *outageLabClient) pollLoop(ctx context.Context) {
	datapage, err := getDataPage(ctx, &olc.options)
	if err != nil {
		fmt.Printf("outagelab: initialization skipped due to error: %v\n", err)
		olc.stop()
		return
	}
	olc.datapage = datapage

	for {
		select {
		case <-time.After(5 * time.Second):
			datapage, err = getDataPage(ctx, &olc.options)
			olc.datapage = datapage
		case <-ctx.Done():
			olc.datapage = nil
			return
		}
	}
}

func getDataPage(ctx context.Context, options *Options) (*datapageResponse, error) {
	reqJson, err := json.Marshal(map[string]string{
		"application": options.Application,
		"environment": options.Environment,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		options.Host+"/api/v1/datapage",
		bytes.NewBuffer(reqJson),
	)
	if err != nil {
		return nil, err
	}

	req.Header["x-api-key"] = []string{options.ApiKey}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == 401 {
		return nil, errors.New("unauthorized request, invalid API key")
	}

	resJson, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var datapage datapageResponse
	err = json.Unmarshal(resJson, &datapage)
	if err != nil {
		return nil, nil
	}

	return &datapage, nil
}

func (olc *outageLabClient) NewTransport(transport http.RoundTripper) *outageLabTransport {
	return newTransport(transport, olc)
}
