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

const (
	intialBackoff = 200 * time.Millisecond
)

type URLStatus struct {
	URL          string
	StatusCode   int
	ResponseTime int64
}

func main() {
	startTime := time.Now()
	urlFileName := flag.String("file", "urls.txt", "File containing URLs (one per line)")
	rateLimitFlag := flag.Bool("rate_lim", false, "Flag to enable rate limiting")
	retries := flag.Int("retries", 0, "Number of retries")
	flag.Parse()

	urls := fetchURLs(*urlFileName)
	urlsStatus := checkURLs(urls, *rateLimitFlag, *retries)
	printURLsStatus(urlsStatus)

	totalTimeTaken := time.Since(startTime)
	fmt.Printf("Execution Time: %v s\n", totalTimeTaken)
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

func checkURLs(urls []string, rateLimitFlag bool, retries int) []URLStatus {
	workerCount := 5

	// create channel for sending jobs
	jobsCh := make(chan string)

	// create channel for receiving results
	resultsCh := make(chan URLStatus, workerCount)

	// create wait group
	var wg sync.WaitGroup
	wg.Add(workerCount)

	// create channel for rate limiting
	rateLimitTicker := time.NewTicker(200 * time.Millisecond)
	defer rateLimitTicker.Stop()
	// time.Tick return <-chan time.Time
	// we will receive a message from this channel in every 200ms (i.e. 5 message every second)

	// create workers
	for i := 0; i < workerCount; i++ {
		go worker(jobsCh, resultsCh, &wg, rateLimitTicker, rateLimitFlag, retries)
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

func worker(
	jobsCh <-chan string, 
	resultsCh chan<- URLStatus, 
	wg *sync.WaitGroup,
	rateLimitTicker *time.Ticker,
	rateLimitFlag bool,
	retries int,
) {
	defer wg.Done()

	httpClient := http.Client{
		Timeout: 5 * time.Second,
	}

	for url := range jobsCh {
		result := checkURL(&httpClient, url, rateLimitTicker, rateLimitFlag, retries)
		resultsCh <- result
	}
}

func checkURL(
	httpClient *http.Client,
	url string, 
	rateLimitTicker *time.Ticker,
	rateLimitFlag bool, 
	retries int, 
) URLStatus {
	backoff := intialBackoff

	for _ = range retries + 1 {
		if rateLimitFlag {
			<- rateLimitTicker.C 
			// Block until a message is received from this channel 
			// this helps us achieve rate limiting 
			// Rate limiting applies per attempt (including retries)
		}

		startTime := time.Now()
		resp, err := httpClient.Get(url)
		responseTime := time.Since(startTime).Milliseconds()
		if err == nil && resp.StatusCode < 500 {
			defer resp.Body.Close()
			return URLStatus{url, resp.StatusCode, responseTime}
		}
		time.Sleep(backoff)
		backoff *= 2
	}
	return URLStatus{url, 0, 0}
}
