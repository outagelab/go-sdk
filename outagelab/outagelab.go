package outagelab

import (
	"fmt"
	"sync"
)

var defaultClient *outageLabClient
var mu sync.Mutex

type Options struct {
	Application string
	Environment string
	ApiKey      string
	Host        string
}

func Start(options Options) {
	if options.ApiKey == "" {
		fmt.Println("outagelab: API key missing, skipping initialization")
		return
	}
	if options.Host == "" {
		options.Host = "https://app.outagelab.com"
	}

	mu.Lock()
	defer mu.Unlock()

	if defaultClient == nil {
		defaultClient = newClient(options)
		defaultClient.start()
	}
}

func Stop() {
	mu.Lock()
	defer mu.Unlock()

	if defaultClient != nil {
		defaultClient.stop()
		defaultClient = nil
	}
}
