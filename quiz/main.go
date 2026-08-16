package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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

	var correctAnswered = 0
	var failedQuestion = 0

	for result, question := range resultQuestionMap {
		fmt.Printf("What is correct answer to: %s \n", question)
		numberOut := getUserInput()

		resultAsNumber, _ := strconv.Atoi(result)

		if numberOut == resultAsNumber {
			correctAnswered++
		} else {
			failedQuestion++
		}
	}

	fmt.Printf("You've completed the quiz with a score of %d failed and %d correct. \n", failedQuestion, correctAnswered)
}
