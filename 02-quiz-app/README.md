# Project 2: Quiz App

This project is a command-line quiz application written in Go.

The application reads quiz questions from a CSV file, presents them to the user one by one, enforces a time limit, and displays the final score at the end.

This project introduces real program structure and multiple core Go concepts working together.

## Requirements
### Functional Requirements
- The quiz questions should be read from a CSV file
- Each question should have:
    - A question
    - A correct answer
- The quiz should:
    - Ask questions one by one
    - Accept user input from standard input
    - Enforce a time limit for the entire quiz
- After the quiz ends (time up or questions finished):
    - Display the total score
    - Display the number of questions attempted
- Example CSV File (problems.csv)
    ```
    5+5,10
    7+3,10
    1+1,2
    8+3,11
    ```
- Example Interaction
    ```bash
    $ go run main.go
    Press Enter to start the quiz...
    Question 1: 5+5 = 10
    Question 2: 7+3 = 10
    Question 3: 1+1 = 2

    Time's up!
    You scored 2 out of 4.
    ```
- Optional Enhancements
- Add command-line flags:
    -limit for quiz duration
    -csv to specify CSV file path
- Shuffle questions
- Show correct answers after quiz ends
- Add per-question timer instead of global timer
