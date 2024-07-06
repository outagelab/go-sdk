package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/outagelab/go-sdk/outagelab"
)

func main() {
	outagelab.Start(outagelab.Options{
		Application: "my-service",
		Environment: "local",
		ApiKey:      os.Getenv("OUTAGELAB_API_KEY"),
		Host:        "http://localhost:8080",
	})
	defer outagelab.Stop()

	for true {
		cnt += 1
		request(cnt, "https://www.google.com")
		request(cnt, "https://vuetifyjs.com")
		time.Sleep(1 * time.Second)
	}
}

var cnt int

func request(cnt int, url string) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("%v: request failed to %v: error: %v\n", cnt, url, err)
	} else {
		fmt.Printf("%v: request to %v: status: %v\n", cnt, url, res.StatusCode)
	}
}
