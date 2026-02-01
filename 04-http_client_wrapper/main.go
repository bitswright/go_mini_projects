package main

import (
	"fmt"
	"net/http"
	"time"
)

type OurClient struct {
	httpClient *http.Client
}

func NewClient() *OurClient {
	return &OurClient{
		httpClient: &http.Client{},
	}
}

func (c *OurClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-CustomWrapperClient", "true")

	start := time.Now()

	resp, err := c.httpClient.Do(req)

	duration := time.Since(start)
	fmt.Println("Request took:", duration)
	return resp, err
}

func main() {
	client := NewClient()

	req, err := http.NewRequest("GET", "https://httpbin.org/get", nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response Status: ", resp.Status)
}
