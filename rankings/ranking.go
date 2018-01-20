package rankings

type Ranking interface {
	Select(int) string
	Insert(int, string)
	Delete(int)
	Length() int
}
