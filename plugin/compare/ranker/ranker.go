package ranker

import (
	"fmt"
	"github.com/mwstobo/rank/input"
	"github.com/mwstobo/rank/rankings"
)

type Ranker struct {
	Ranking   rankings.Ranking
	userInput input.UserInput
}

func NewRanker(ranking rankings.Ranking, userInput input.UserInput) *Ranker {
	return &Ranker{
		Ranking:   ranking,
		userInput: userInput,
	}
}

func (ranker *Ranker) AddItem(item string) error {
	var baseIndex, comparativeIndex, rankingLength, middle int
	var comparativeItem string

	rankingLength = ranker.Ranking.Length()
	baseIndex = 0

	for {
		middle = rankingLength / 2
		comparativeIndex = baseIndex + middle
		if rankingLength == 0 {
			break
		}

		comparativeItem = ranker.Ranking.Select(comparativeIndex)

		isHigher := func() {
			rankingLength -= rankingLength - middle
		}

		isLower := func() {
			rankingLength -= middle + 1
			baseIndex += middle + 1
		}

		choices := []input.Choice{
			input.Choice{Command: "y", Action: isHigher},
			input.Choice{Command: "n", Action: isLower},
		}

		selectedChoice, err := ranker.userInput.PresentChoice(
			fmt.Sprintf("Is %s higher than %s?", item, comparativeItem),
			choices)
		if err != nil {
			return err
		}
		selectedChoice()
	}

	ranker.Ranking.Insert(comparativeIndex, item)
	return nil
}

func (ranker *Ranker) ListItems() {
	if ranker.Ranking.Length() == 0 {
		fmt.Println("No items")
		return
	}

	for i := 0; i < ranker.Ranking.Length(); i += 1 {
		item := ranker.Ranking.Select(i)
		fmt.Printf("%d. %s\n", i+1, item)
	}
}

func (ranker *Ranker) DeleteItem(itemIndex int) {
	ranker.Ranking.Delete(itemIndex)
}
