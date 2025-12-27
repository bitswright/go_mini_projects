# Project 1: Number Guessing Game

This is a command-line number guessing game written in Go.

The program generates a random number within a fixed range, and the user must guess the number. After each guess, the program provides feedback until the correct number is guessed.

This project introduces core Go programming concepts such as variables, loops, user input, randomness, and basic error handling.

## Requirements
### Functional Requirements
- The program should generate a random number between 1 and 100.
- The user should repeatedly enter guesses via standard input.
- After each guess: 
    - Inform the user if the guess is too high or too low.
- When the user guesses correctly: 
    - Print a success message
    - Display the number of attempts taken
    - Exit the program
- Example Usage
    ```bash
    $ go run main.go
    Hello, World!
    
    $ go run main.go Harish
    Hello, Harish!
    ```