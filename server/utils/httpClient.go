package utils

import "net/http"

var httpClient = &http.Client{}

func HttpClient() *http.Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return httpClient
}
