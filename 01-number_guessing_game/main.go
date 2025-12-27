package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("Welcome to the Guessing Game!")

	rand.Seed(time.Now().UnixNano())
	randomNumberSelected := rand.Intn(100) + 1
	fmt.Println("I have picked a number between 1 and 100.")

	var numberGuessed int
	attempts := 0
	for {
		fmt.Printf("Enter your guess: ")
		_, err := fmt.Scanf("%d", &numberGuessed)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		attempts++
		if numberGuessed == randomNumberSelected {
			fmt.Println("Correct!")
			fmt.Printf("You guessed the number in %d attempts.\n", attempts)
			return
		} else if numberGuessed < randomNumberSelected {
			if randomNumberSelected - numberGuessed < 10 {
				fmt.Println("Bit low!")
			} else {
				fmt.Println("Too low!")
			}
		} else {
			if numberGuessed - randomNumberSelected < 10 {
				fmt.Println("Bit high!")
			} else {
				fmt.Println("Too high!")
			}
		}
	}
}