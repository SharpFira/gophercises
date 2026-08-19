package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

func getFile(path string) *os.File {
	file, error := os.Open(path)

	if error != nil {
		errorMessage := fmt.Sprintf("Could not read file: %s", path)
		log.Fatal(errorMessage)
	}

	return file
}

func getUserInput() int {
	reader := bufio.NewReader(os.Stdin)
	line, error := reader.ReadString('\n')

	if error != nil {
		log.Fatal("Could not process user input")
	}

	var numberIn, inputError = strconv.Atoi(strings.TrimSpace(line))

	if inputError != nil {
		fmt.Println("Wrong input. Type a number")
		getUserInput()
	}

	return numberIn
}

func getTimerFlag() int {
	fmt.Println("Ready to start the quiz. The default time is set to 10 seconds but you can modify it.")

	fmt.Println("Do you wish to modify timer? write yes/no")

	userWantsToCustomize := false
	defaultTimerNumber := 30
	var attempt, maxAttempts int = 0, 3
	var possibleAnswers = []string{"yes", "no"}

	for attempt <= maxAttempts {

		reader := bufio.NewReader(os.Stdin)

		line, error := reader.ReadString('\n')

		if error != nil {
			log.Fatal("User has stopped or done something wrong...")

			os.Exit(3)
		}

		userInput := strings.TrimSpace(strings.ToLower(line))

		if !slices.Contains(possibleAnswers, userInput) {
			attempt++
			fmt.Printf("Wrong answer - attempt=%d. Try answering again?", attempt)
			fmt.Print("\n")
		} else {

			if userInput == possibleAnswers[1] {
				break
			}

			userWantsToCustomize = true

			break
		}
	}

	var timerInput int = defaultTimerNumber

	if userWantsToCustomize {
		fmt.Print("Enter a number in seconds you want the timer to run: ")
		fmt.Scan(&timerInput)
		fmt.Println("The timer will run for: ", timerInput)
	}

	return timerInput
}

func main() {
	csvPath := "./problems.csv"

	file := getFile(csvPath)

	defer file.Close()

	reader := csv.NewReader(file)

	records, error := reader.ReadAll()

	if error != nil {
		errorMessage := "Error reading rows from csv."
		log.Fatal(errorMessage)
	}

	resultQuestionMap := make(map[string]string)

	for i := 0; i < len(records); i++ {

		question := records[i][0]
		result := records[i][1]

		resultQuestionMap[result] = question
	}

	lengthOfMap := strconv.Itoa(len(resultQuestionMap))
	fmt.Printf("The csv has been read and there are %s in memory. \n", lengthOfMap)

	var selectedTimer int = getTimerFlag()

	fmt.Println("The user have set the timer to: ", selectedTimer)

	var correctAnswered = 0
	var failedQuestion = 0

	timer := time.NewTimer(time.Duration(selectedTimer) * time.Second)
	defer timer.Stop()

	answerCh := make(chan int, 1)
	expired := false

	for result, question := range resultQuestionMap {
		if expired {
			break
		}

		fmt.Printf("What is correct answer to: %s \n", question)

		go func() {
			answerCh <- getUserInput()
		}()

		select {
		case <-timer.C:
			fmt.Println("\nTime's up!")
			expired = true

		case numberOut := <-answerCh:
			resultAsNumber, _ := strconv.Atoi(result)

			if numberOut == resultAsNumber {
				correctAnswered++
			} else {
				failedQuestion++
			}
		}
	}

	fmt.Printf("You've completed the quiz with a score of %d failed and %d correct. \n", failedQuestion, correctAnswered)
}
