package input

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type StdinUserInput struct{}

func NewStdinUserInput() UserInput {
	return &StdinUserInput{}
}

func (userInput *StdinUserInput) VerbosePresentChoice(
	prompt string,
	choices []Choice) (func(), error) {

	for _, choice := range choices {
		fmt.Printf("%s) %s\n", choice.Command, choice.Name)
	}

	return userInput.PresentChoice(prompt, choices)
}
func (userInput *StdinUserInput) PresentChoice(
	prompt string,
	choices []Choice) (func(), error) {

	choiceMap := make(map[string]func())

	choiceCommands := []string{}
	for _, choice := range choices {
		choiceCommands = append(choiceCommands, choice.Command)
		choiceMap[choice.Command] = choice.Action
	}

	stdinScanner := bufio.NewScanner(os.Stdin)

	for {
		choiceString := strings.Join(choiceCommands, "/")
		fmt.Printf("%s: [%s] ", prompt, choiceString)

		moreTokens := stdinScanner.Scan()
		if !moreTokens {
			fmt.Println("")
			if err := stdinScanner.Err(); err != nil {
				return nil, err
			}
			return nil, EOF
		}
		choice := stdinScanner.Text()

		if choiceFunction, ok := choiceMap[choice]; ok {
			return choiceFunction, nil
		} else {
			fmt.Printf("Invalid choice %s\n", choice)
		}
	}
}

func (userInput *StdinUserInput) GetInput(prompt string) (string, error) {
	stdinScanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("%s ", prompt)

		moreTokens := stdinScanner.Scan()
		if !moreTokens {
			fmt.Println("")
			if err := stdinScanner.Err(); err != nil {
				return "", err
			}
			return "", EOF
		}

		return stdinScanner.Text(), nil
	}
}

func (userInput *StdinUserInput) GetInputInt(prompt string) (int, error) {
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
