package ranker_test

import (
	"github.com/mwstobo/rank/mocks"
	. "github.com/mwstobo/rank/ranker"
	"github.com/mwstobo/rank/rankings"
	"testing"
)

func TestAddItem(t *testing.T) {
	t.Run("First", AddFirst)
	t.Run("Best", AddBest)
	t.Run("Worst", AddWorst)
	t.Run("MiddleOdd", AddMiddleOdd)
	t.Run("MiddleEven", AddMiddleEven)
}

func AddFirst(t *testing.T) {
	rankingSlice := []string{}
	ranking := rankings.NewArrayRanking(rankingSlice)

	mockUserInput := mocks.NewMockUserInput([]string{}, []string{})

	ranker := NewRanker(ranking, mockUserInput)
	ranker.AddItem("itemA")

	expected := "itemA"
	actual := ranker.Ranking.Select(0)

	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func AddBest(t *testing.T) {
	rankingSlice := []string{"itemA"}
	ranking := rankings.NewArrayRanking(rankingSlice)

	mockUserInput := mocks.NewMockUserInput([]string{"y"}, []string{})

	ranker := NewRanker(ranking, mockUserInput)
	ranker.AddItem("itemB")

	expected := "itemB"
	actual := ranker.Ranking.Select(0)

	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func AddWorst(t *testing.T) {
	rankingSlice := []string{"itemA"}
	ranking := rankings.NewArrayRanking(rankingSlice)

	mockUserInput := mocks.NewMockUserInput([]string{"n"}, []string{})

	ranker := NewRanker(ranking, mockUserInput)
	ranker.AddItem("itemB")

	expected := "itemB"
	actual := ranker.Ranking.Select(ranker.Ranking.Length() - 1)

	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func AddMiddleOdd(t *testing.T) {
	rankingSlice := []string{"itemA", "itemB", "itemC"}
	ranking := rankings.NewArrayRanking(rankingSlice)

	mockUserInput := mocks.NewMockUserInput([]string{"y", "n"}, []string{})

	ranker := NewRanker(ranking, mockUserInput)
	ranker.AddItem("itemD")

	expected := "itemD"
	actual := ranker.Ranking.Select(1)

	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func AddMiddleEven(t *testing.T) {
	rankingSlice := []string{"itemA", "itemB"}
	ranking := rankings.NewArrayRanking(rankingSlice)

	mockUserInput := mocks.NewMockUserInput([]string{"y", "n"}, []string{})

	ranker := NewRanker(ranking, mockUserInput)
	ranker.AddItem("itemC")

	expected := "itemC"
	actual := ranker.Ranking.Select(1)

	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
