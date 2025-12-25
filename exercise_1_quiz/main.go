package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type problem struct {
	question, answer string
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "A CSV file in the format of 'question,answer'")
	timeLimit := flag.Int("timeLimit", 30, "Time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %s", *csvFilename)
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Failed to parse the provided CSV file")
	}

	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0
	for i, currProblem := range problems {
		fmt.Printf("Problem #%d: %s is ", i+1, currProblem.question)
		answerChan := make(chan string)
		// This is a anonymous function also a goroutine
		// This is also a closure (using variables declared in the outer function)
		go func() {
			var answer string
			fmt.Scanf("%s", &answer)
			answerChan <- answer
		}()
		select {
		case <-timer.C:
			fmt.Printf("\nYou have scored %d out of %d\n", correct, len(problems))
			return
		case answerGiven := <-answerChan:
			if answerGiven == currProblem.answer {
				correct++
			}
		}
	}
	fmt.Printf("You have scored %d out of %d\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))

	for i, line := range lines {
		ret[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}

	return ret
}
