package ui

import (
	"fmt"
	"github.com/mwstobo/rank/input"
	"github.com/mwstobo/rank/ranker"
	"github.com/mwstobo/rank/store"
)

type InteractiveApp struct {
	ranker    *ranker.Ranker
	storage   store.Storage
	userInput input.UserInput
	quit      chan bool
}

func NewInteractiveApp(
	ranker *ranker.Ranker,
	storage store.Storage,
	userInput input.UserInput) *InteractiveApp {

	return &InteractiveApp{
		ranker:    ranker,
		storage:   storage,
		userInput: userInput,
		quit:      make(chan bool, 1),
	}
}

func (app *InteractiveApp) Run() {
	choices := []input.Choice{
		input.Choice{"Add an item", "1", app.AddAction},
		input.Choice{"View your list", "2", app.ListAction},
		input.Choice{"Delete an item", "3", app.DeleteAction},
		input.Choice{"Save and quit", "4", app.SaveAndQuitAction},
		input.Choice{"Quit without saving", "5", app.QuitAction},
	}

	for {
		select {
		case <-app.quit:
			return
		default:
			selectedChoice, err := app.userInput.VerbosePresentChoice(
				"What do you want to do?",
				choices)
			if err == input.EOF {
				app.SaveAndQuitAction()
				continue
			} else if err != nil {
				fmt.Printf("Error reading choice: %v\n", err)
				app.SaveAndQuitAction()
				continue
			}
			selectedChoice()
			fmt.Println("")
		}
	}
}

func (app *InteractiveApp) AddAction() {
	item, err := app.userInput.GetInput("What do you want to add?")
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		app.SaveAndQuitAction()
		return
	}

	err = app.ranker.AddItem(item)
	if err != nil {
		fmt.Printf("Error adding item: %v\n", err)
		app.SaveAndQuitAction()
		return
	}
}

func (app *InteractiveApp) ListAction() {
	app.ranker.ListItems()
}

func (app *InteractiveApp) DeleteAction() {
	var itemNumber int
	var err error

	if app.ranker.Ranking.Length() == 0 {
		fmt.Println("No items")
		return
	}

	app.ListAction()
	for {
		itemNumber, err = app.userInput.GetInputInt(
			"What do you want to delete?")
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			app.SaveAndQuitAction()
			return
		}

		if itemNumber > 0 && itemNumber <= app.ranker.Ranking.Length() {
			break
		}
		fmt.Printf("Invalid item number %d\n", itemNumber)
	}

	app.ranker.DeleteItem(itemNumber - 1)
}

func (app *InteractiveApp) SaveAndQuitAction() {
	app.storage.Export(app.ranker.Ranking)
	app.QuitAction()
}

func (app *InteractiveApp) QuitAction() {
	app.quit <- true
}
