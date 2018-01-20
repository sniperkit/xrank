package input

import (
	"errors"
)

var EOF = errors.New("user closed stdin")

type Choice struct {
	Name    string
	Command string
	Action  func()
}

type UserInput interface {
	VerbosePresentChoice(string, []Choice) (func(), error)
	PresentChoice(string, []Choice) (func(), error)
	GetInput(string) (string, error)
	GetInputInt(string) (int, error)
}
