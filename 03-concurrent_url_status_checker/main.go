package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"text/tabwriter"
	"time"
)

type URLStatus struct {
	URL          string
	StatusCode   int
	ResponseTime int64
}

func main() {
	urlFileName := flag.String("file", "urls.txt", "File containing URLs (one per line)")
	flag.Parse()

	urls := fetchURLs(*urlFileName)
	urlsStatus := checkURLs(urls)
	printURLsStatus(urlsStatus)
}

func fetchURLs(fileName string) (urls []string) {
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

func checkURLs(urls []string) []URLStatus {
	workerCount := 5

	// create channel for sending jobs
	jobsCh := make(chan string)

	// create channel for receiving results
	resultsCh := make(chan URLStatus, workerCount)

	// create wait group
	var wg sync.WaitGroup
	wg.Add(workerCount)

	// create workers
	for i := 0; i < workerCount; i++ {
		go worker(jobsCh, resultsCh, &wg)
	}

	// send urls on jobCh and close the channel
	// as no more urls are to be sent for checking by workers
	go func() {
		for _, url := range urls {
			jobsCh <- url
		}
		close(jobsCh)
	}()

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	// receive results from resultsCh
	var urlsStatus []URLStatus
	for result := range resultsCh {
		urlsStatus = append(urlsStatus, result)
	}

	return urlsStatus
}

func printURLsStatus(urlsStatus []URLStatus) {
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

func worker(jobsCh <-chan string, resultsCh chan<- URLStatus, wg *sync.WaitGroup) {
	defer wg.Done()

	httpClient := http.Client{
		Timeout: 5 * time.Second,
	}

	for url := range jobsCh {
		startTime := time.Now()
		resp, err := httpClient.Get(url)
		responseTime := time.Since(startTime).Milliseconds()
		if err != nil {
			resultsCh <- URLStatus{url, 0, 0}
			continue
		}
		resultsCh <- URLStatus{url, resp.StatusCode, responseTime}
		resp.Body.Close()
	}
}
