package common

import (
	"fmt"
	"net/http"
	"time"
)

const USER_AGENT = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

var defaultClient = &http.Client{Timeout: 60 * time.Second}

type StockNav interface {
	GetNav() float64
	GetNavDate() time.Time
}

func GetHttpClient() *http.Client {
	return defaultClient
}

func FireRequest(r *http.Request) (*http.Response, error) {
	resp, err := GetHttpClient().Do(r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("Received status code: %d, instead of 200", resp.StatusCode)
	}
	return resp, nil
}
