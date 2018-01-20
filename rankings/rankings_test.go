package rankings_test

import (
	. "github.com/mwstobo/rank/rankings"
	"testing"
)

func TestRankings(t *testing.T) {
	rankingConstructors := map[string]func([]string) Ranking{
		"array": NewArrayRanking,
	}

	for name, rankingConstructor := range rankingConstructors {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			SelectSuite(t, rankingConstructor)
			InsertSuite(t, rankingConstructor)
			DeleteSuite(t, rankingConstructor)
			LengthSuite(t, rankingConstructor)
		})
	}
}

func SelectSuite(t *testing.T, rankingConstructor func([]string) Ranking) {
	rankingSlice := []string{
		"item0",
		"item1",
		"item2",
	}
	arrayRanking := rankingConstructor(rankingSlice)

	t.Run("FirstItem", func(t *testing.T) {
		SelectFirstItem(t, arrayRanking, rankingSlice)
	})
	t.Run("LastItem", func(t *testing.T) {
		SelectLastItem(t, arrayRanking, rankingSlice)
	})
	t.Run("MiddleItem", func(t *testing.T) {
		SelectMiddleItem(t, arrayRanking, rankingSlice)
	})
}

func SelectFirstItem(
	t *testing.T,
	arrayRanking Ranking,
	rankingSlice []string) {

	expected := rankingSlice[0]
	actual := arrayRanking.Select(0)

	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func SelectLastItem(
	t *testing.T,
	arrayRanking Ranking,
	rankingSlice []string) {

	expected := rankingSlice[2]
	actual := arrayRanking.Select(2)

	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func SelectMiddleItem(
	t *testing.T,
	arrayRanking Ranking,
	rankingSlice []string) {

	expected := rankingSlice[1]
	actual := arrayRanking.Select(1)

	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func InsertSuite(t *testing.T, rankingConstructor func([]string) Ranking) {
	t.Run("FirstItem", func(t *testing.T) {
		rankingSlice := []string{"item"}
		arrayRanking := rankingConstructor(rankingSlice)
		InsertFirstItem(t, arrayRanking, rankingSlice)
	})
	t.Run("LastItem", func(t *testing.T) {
		rankingSlice := []string{"item"}
		arrayRanking := rankingConstructor(rankingSlice)
		InsertLastItem(t, arrayRanking, rankingSlice)
	})
}

func InsertFirstItem(
	t *testing.T,
	arrayRanking Ranking,
	rankingSlice []string) {

	arrayRanking.Insert(0, "itemFirst")

	expected := "itemFirst"
	actual := arrayRanking.Select(0)

	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func InsertLastItem(
	t *testing.T,
	arrayRanking Ranking,
	rankingSlice []string) {

	arrayRanking.Insert(len(rankingSlice), "itemLast")

	expected := "itemLast"
	actual := arrayRanking.Select(len(rankingSlice))

	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func DeleteSuite(t *testing.T, rankingConstructor func([]string) Ranking) {
	t.Run("FirstItem", func(t *testing.T) {
		rankingSlice := []string{
			"item0",
			"item1",
		}
		arrayRanking := rankingConstructor(rankingSlice)
		DeleteFirstItem(t, arrayRanking, rankingSlice)
	})
	t.Run("LastItem", func(t *testing.T) {
		rankingSlice := []string{
			"item0",
			"item1",
		}
		arrayRanking := rankingConstructor(rankingSlice)
		DeleteLastItem(t, arrayRanking, rankingSlice)
	})
}

func DeleteFirstItem(
	t *testing.T,
	arrayRanking Ranking,
	rankingSlice []string) {

	initialLength := len(rankingSlice)

	expected := arrayRanking.Select(1)
	arrayRanking.Delete(0)
	actual := arrayRanking.Select(0)

	expectedLength := initialLength - 1
	actualLength := arrayRanking.Length()

	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}

	if expectedLength != actualLength {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func DeleteLastItem(
	t *testing.T,
	arrayRanking Ranking,
	rankingSlice []string) {

	initialLength := len(rankingSlice)

	expected := arrayRanking.Select(0)
	arrayRanking.Delete(1)
	actual := arrayRanking.Select(0)

	expectedLength := initialLength - 1
	actualLength := arrayRanking.Length()

	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}

	if expectedLength != actualLength {
		t.Errorf("Expected %d, got %d", expectedLength, actualLength)
	}
}

func LengthSuite(t *testing.T, rankingConstructor func([]string) Ranking) {
	t.Run("Zero", func(t *testing.T) {
		rankingSlice := []string{}
		arrayRanking := rankingConstructor(rankingSlice)
		Length(t, arrayRanking, rankingSlice)
	})
	t.Run("NonZero", func(t *testing.T) {
		rankingSlice := []string{
			"item0",
		}
		arrayRanking := rankingConstructor(rankingSlice)
		Length(t, arrayRanking, rankingSlice)
	})
}

func Length(
	t *testing.T,
	arrayRanking Ranking,
	rankingSlice []string) {

	expected := len(rankingSlice)
	actual := arrayRanking.Length()

	if expected != actual {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}
