# Exercise 1: Quiz Game

## Requirements
### Part-1
- Create a quiz game, that reads a CSV file and give the questions from the file to the user keeping tracking of the questions they have attempted and how many they got right.
- Irrespective of correctness of the user's answer, next question should be asked immediately afterwards.
- CSV file should default to `problems.csv`, but the user could customize the file using a flag.
- CSV file should consist of two columns, question and respective answers.
- Expect the answers to be one word/number.
- Expect < 100 questions in the CSV file.
- At the end of the quiz, program should output total number of questions and total number of correct answers.

### Part-2
- Add a timer for the quiz, with default value as 30 seconds, but customizable via a flag.
- End the quiz immediately after the timer exceeds, even if answers are expected from user.
- User should be asked to press `Enter` (or some other key) before starting the timer and then only the questions should be printed out on the screen to the user.

### Bonus
- Add string trimming and cleanup to ensure correct answers with extra whitespace, capitalization etc are not considered incorrect.
- Add an option (a new flag) to shuffle the quiz each time it is run.
