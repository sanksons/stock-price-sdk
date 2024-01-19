package bse

import (
	"net/http"
	"time"
)

var defaultClient = &http.Client{Timeout: 60 * time.Second}

func GetHttpClient() *http.Client {
	return defaultClient
}
