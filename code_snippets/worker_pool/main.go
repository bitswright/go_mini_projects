package main

import (
	"fmt"
	"time"
)

func main() {
	// create jobsCh channel: channel on which jobs will be given to the workers
	jobsCh := make(chan string)
	defer close(jobsCh)

	// create workers
	go worker(1, jobsCh)
	go worker(2, jobsCh)
	go worker(3, jobsCh)

	// send jobs on jobsCh
	jobsCh <- "A"
	jobsCh <- "B"
	jobsCh <- "C"
	jobsCh <- "D"
	jobsCh <- "E"
	jobsCh <- "F"

	time.Sleep(2*time.Second)
}

func worker(id int, job <-chan string) {
	for text := range job {
		fmt.Println("Worker %d performing job %s", id, text)
	}
}