package outagelab

type datapageResponse struct {
	Rules []*outageRule `json:"rules"`
}

type outageRule struct {
	Type                string               `json:"type"`
	HttpClientRequestV1 *httpClientRequestV1 `json:"httpClientRequestV1"`
}

type httpClientRequestV1 struct {
	Host     string `json:"host"`
	Status   int    `json:"status"`
	Duration int    `json:"Duration"`
}
