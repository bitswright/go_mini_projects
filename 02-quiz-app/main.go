package main

import (
	"fmt"
	"os"
	"log"
	"encoding/csv"
)

type problem struct {
	question 	string
	answer 		string
}

func main() {
	problemsFileName := "problems.csv"

	problems := getProblems(problemsFileName)
	
	score := take_test(problems)

	fmt.Printf("You score %d out of %d.\n", score, len(problems))
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

func take_test(problems []problem) int {
	correctAnswers := 0
	for i, problem := range problems {
		fmt.Printf("Problem#%d: %s is ", i, problem.question)
		var answerByUser string
		fmt.Scanf("%s", &answerByUser)
		if answerByUser == problem.answer {
			correctAnswers++;
		}
	}
	return correctAnswers;
}