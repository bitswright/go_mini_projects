package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

type problem struct {
	question string
	answer   string
}

func main() {
	problemsFileName := flag.String("csv_file_name", "problems.csv", "Path to the CSV file containing the quiz problems")
	timeLimitInSeconds := flag.Int("quiz_time_limit", 30, "Time limit for the quiz in seconds")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	problems := getProblems(*problemsFileName)

	shuffleProblems(problems)

	score, questionsAttempted := takeTest(problems, *timeLimitInSeconds)

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

func takeTest(problems []problem, timeLimitInSeconds int) (int, int) {
	correctAnswers := 0
	var answerByUser string
	var quizTimerCh <-chan time.Time

	quizTimerCh = startTimer(timeLimitInSeconds)

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
		case <-quizTimerCh:
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

func startTimer(timeLimitInSeconds int) <-chan time.Time {
	return time.After(time.Duration(timeLimitInSeconds) * time.Second)
}

func shuffleProblems(problems []problem) {
	rand.Shuffle(len(problems), func(i, j int) {
		problems[i], problems[j] = problems[j], problems[i]
	})
}
