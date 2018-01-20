// +build mocks

package mocks

import (
	. "github.com/mwstobo/rank/input"
	"strconv"
)

type MockUserInput struct {
	choiceChan chan string
	inputChan  chan string
}

func NewMockUserInput(
	choiceInputs []string,
	normalInputs []string) UserInput {

	userInput := &MockUserInput{
		choiceChan: make(chan string),
		inputChan:  make(chan string),
	}

	go userInput.generateInput(choiceInputs, userInput.choiceChan)
	go userInput.generateInput(normalInputs, userInput.inputChan)

	return userInput
}

func (userInput *MockUserInput) VerbosePresentChoice(
	prompt string,
	choices []Choice) (func(), error) {

	return userInput.PresentChoice(prompt, choices)
}
func (userInput *MockUserInput) PresentChoice(
	prompt string,
	choices []Choice) (func(), error) {

	choiceMap := make(map[string]func())

	for _, choice := range choices {
		choiceMap[choice.Command] = choice.Action
	}

	return choiceMap[<-userInput.choiceChan], nil
}

func (userInput *MockUserInput) GetInput(prompt string) (string, error) {
	return <-userInput.inputChan, nil
}

func (userInput *MockUserInput) GetInputInt(prompt string) (int, error) {
	text, err := userInput.GetInput(prompt)
	if err != nil {
		return 0, err
	}

	integer, err := strconv.Atoi(text)
	if err != nil {
		return 0, err
	}

	return integer, nil
}

func (userInput *MockUserInput) generateInput(
	inputs []string,
	inputChan chan string) {

	position := 0
	for {
		if len(inputs) == 0 {
			continue
		}
		inputChan <- inputs[position%len(inputs)]
		position += 1
	}
}
