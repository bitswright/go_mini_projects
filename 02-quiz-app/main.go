package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

type problem struct {
	question string
	answer   string
}

func main() {
	problemsFileName := "problems.csv"
	timeLimitInSeconds := 2

	problems := getProblems(problemsFileName)

	score, questionsAttempted := take_test(problems, timeLimitInSeconds)

	fmt.Printf("You have attempted %d question(s). You score %d out of %d.\n", questionsAttempted, score, len(problems))
}

func getAllRecordsFromCSVFile(fileName string) [][]string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error while reading the file %s. Error: %s\n", fileName, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error while reading CSV file %s.\n", fileName)
	}

	return records
}

func getProblems(problemsFileName string) []problem {
	records := getAllRecordsFromCSVFile(problemsFileName)

	problems := make([]problem, len(records))
	for i, record := range records {
		problems[i] = problem{record[0], record[1]}
	}

	return problems
}

func take_test(problems []problem, timeLimitInSeconds int) (int, int) {
	correctAnswers := 0
	var answerByUser string

	timer := time.NewTimer(time.Duration(timeLimitInSeconds) * time.Second)

	for i, problem := range problems {
		fmt.Printf("Problem#%d: %s is ", i, problem.question)

		// take user input
		inputChannel := make(chan string)
		go func() {
			var input string
			fmt.Scanln(&input)
			inputChannel <- input
		}()

		// wait for user input or timeout
		select {
		case <-timer.C:
			fmt.Println("\nTime is up!")
			return correctAnswers, i
		case answerByUser = <-inputChannel:
			if answerByUser == problem.answer {
				correctAnswers++
			}
		}
	}
	return correctAnswers, len(problems)
}
