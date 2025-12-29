package main

import (
	"flag"
	"os"
	"log"
	"bufio"
	"strings"
	"fmt"
	"time"
	"net/http"
	"text/tabwriter"
)

type URLStatus struct {
	URL          string
	StatusCode   int
	ResponseTime int64
}

func main() {
	URLFileName := flag.String("file", "urls.txt", "File containing URLs (one per line)")
	flag.Parse()

	urls := fetchUrls(*URLFileName)

	urlsStatus := checkUrls(urls)
	printUrlsStatus(urlsStatus)
}

func fetchUrls(fileName string) (urls []string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error while reading the file %s, error: %s", fileName, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url != "" {
			fullURL := getFullURL(url)
			urls = append(urls, fullURL)
		}
	}
	if err = scanner.Err(); err != nil {
		log.Fatalf("Error while reading the file %s, error: %s", fileName, err)
	}

	return
}

func checkUrls(urls []string) []URLStatus {
	urlsStatus := make([]URLStatus, len(urls))
	httpClient := http.Client{
		Timeout: 5 * time.Second,
	}
	
	for i, url := range urls {
		startTime := time.Now()
		resp, err := httpClient.Get(url)
		responseTime := time.Since(startTime).Milliseconds()
		if err != nil {
			urlsStatus[i] = URLStatus{url, 0, 0} 
		} else {
			urlsStatus[i] = URLStatus{url, resp.StatusCode, responseTime} 
		}
		if resp != nil {
			resp.Body.Close()
		}
	}

	return urlsStatus
}

func printUrlsStatus(urlsStatus []URLStatus) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(writer, "URL\tSTATUS\tRESPONSE TIME")
	for _, URLStatus := range urlsStatus {
		var statusCodeStr string
		var responseTimeStr string
		if URLStatus.StatusCode == 0 {
			statusCodeStr = "ERROR"
			responseTimeStr = "timeout"
		} else {
			statusCodeStr = fmt.Sprintf("%d", URLStatus.StatusCode)
			responseTimeStr = fmt.Sprintf("%d ms", URLStatus.ResponseTime)
		}
		fmt.Fprintf(
			writer,
			"%s\t%s\t%s\n",
			URLStatus.URL,
			statusCodeStr,
			responseTimeStr,
		)
	}
	writer.Flush()
}

func getFullURL(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}
	return url
}