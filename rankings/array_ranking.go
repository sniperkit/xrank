package rankings

type ArrayRanking struct {
	array []string
}

func NewArrayRanking(rankingSlice []string) Ranking {
	return &ArrayRanking{
		array: rankingSlice,
	}
}

func (ranking *ArrayRanking) Select(index int) string {
	return ranking.array[index]
}

func (ranking *ArrayRanking) Insert(index int, item string) {
	ranking.array = append(ranking.array, "")
	copy(ranking.array[index+1:], ranking.array[index:])
	ranking.array[index] = item
}

func (ranking *ArrayRanking) Delete(index int) {
	ranking.array = append(ranking.array[:index], ranking.array[index+1:]...)
}

func (ranking *ArrayRanking) Length() int {
	return len(ranking.array)
}
