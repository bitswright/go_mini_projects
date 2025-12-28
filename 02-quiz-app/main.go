package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"text/tabwriter"
	"time"
)

type problem struct {
	question string
	answer   string
}

type problemAttempt struct {
	problem         problem
	answerAttempted string
}

func main() {
	problemsFileName := flag.String("csv_file_name", "problems.csv", "Path to the CSV file containing the quiz problems")
	timeLimitInSeconds := flag.Int("quiz_time_limit", 30, "Time limit for the quiz in seconds")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	problems := getProblems(*problemsFileName)
	shuffleProblems(problems)
	score, questionsAttempted, incorrectProblemsAttempts := takeTest(problems, *timeLimitInSeconds)
	printScoreAndCorrectAnswers(score, questionsAttempted, len(problems), incorrectProblemsAttempts)
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

func takeTest(problems []problem, timeLimitInSeconds int) (int, int, []problemAttempt) {
	fmt.Printf("Press Enter to start test...")
	bufio.NewReader(os.Stdin).ReadString('\n')

	correctAnswers := 0
	var incorrectProblemsAttempted []problemAttempt
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
			return correctAnswers, i, incorrectProblemsAttempted
		case answerByUser = <-inputChannel:
			if answerByUser == problem.answer {
				correctAnswers++
			} else {
				incorrectProblemsAttempted = append(incorrectProblemsAttempted, problemAttempt{problem, answerByUser})
			}
		}
	}
	return correctAnswers, len(problems), incorrectProblemsAttempted
}

func startTimer(timeLimitInSeconds int) <-chan time.Time {
	return time.After(time.Duration(timeLimitInSeconds) * time.Second)
}

func shuffleProblems(problems []problem) {
	rand.Shuffle(len(problems), func(i, j int) {
		problems[i], problems[j] = problems[j], problems[i]
	})
}

func printScoreAndCorrectAnswers(score, questionsAttempted, problemsCount int, incorrectProblemAttempts []problemAttempt) {
	fmt.Printf("You have attempted %d question(s). You score %d out of %d.\n", questionsAttempted, score, problemsCount)
	if len(incorrectProblemAttempts) > 0 {
		fmt.Println("Incorrect answers solution")
		fmt.Println("==========================")
		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(writer, "Question\tAttempted Answer\tCorrect Answer")
		for _, p := range incorrectProblemAttempts {
			fmt.Fprintf(
				writer,
				"%s\t%s\t%s\n",
				p.problem.question,
				p.answerAttempted,
				p.problem.answer,
			)
		}
		writer.Flush()
	}
}
