package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
)

const (
	INTERACTIVE = "interactive"
	CLI         = "cli"

	ADD_COMMAND    = "add"
	LIST_COMMAND   = "list"
	DELETE_COMMAND = "delete"
)

var (
	Ui               string
	Filename         string
	Command          string
	AddItem          string
	DeleteItemNumber int

	Help = errors.New("user asked for help")
)

func ParseConfig() error {
	commandActions := buildCommandActions()
	Ui = INTERACTIVE

	args := os.Args[1:]

	if len(args) <= 0 || args[0] == "--help" || args[0] == "-h" {
		Usage()
		return Help
	}

	switch len(args) {
	case 0:
		return errors.New("no ranking filename")
	case 1:
		Filename = args[0]
	default:
		Filename = args[0]
		Command = args[1]
		Ui = CLI
		args = args[2:]

		if commandAction, ok := commandActions[Command]; !ok {
			return fmt.Errorf("invalid command %s", Command)
		} else {
			err := commandAction(args)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Usage() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintln(w, "Usage: rank FILENAME [COMMAND]")
	fmt.Fprintln(w,
		"If you do not provide COMMAND, rank will start in interactive mode")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "COMMANDs")
	fmt.Fprintln(w, "\tadd ITEM\tadd ITEM to ranking")
	fmt.Fprintln(w, "\tlist\tlist all of the items in ranking")
	fmt.Fprintln(w, "\tdelete ITEM_NUMBER\tdelete ITEM_NUMBER from ranking")
	w.Flush()
}

func buildCommandActions() map[string]func([]string) error {
	return map[string]func([]string) error{
		ADD_COMMAND:    parseAddConfig,
		LIST_COMMAND:   parseListConfig,
		DELETE_COMMAND: parseDeleteConfig,
	}
}

func parseAddConfig(args []string) error {
	switch len(args) {
	case 1:
		AddItem = args[0]
	default:
		return errors.New("missing add item")
	}
	return nil
}

func parseListConfig(args []string) error {
	switch len(args) {
	case 0:
		return nil
	default:
		return errors.New("too many args for list")
	}
}

func parseDeleteConfig(args []string) error {
	switch len(args) {
	case 1:
		itemNumberString := args[0]
		if itemNumber, err := strconv.Atoi(itemNumberString); err != nil {
			return fmt.Errorf("can't convert %s to int: %v", err)
		} else {
			DeleteItemNumber = itemNumber
		}
	default:
		return errors.New("missing delete item number")
	}
	return nil
}
