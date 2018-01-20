package store

import (
	"github.com/mwstobo/rank/rankings"
)

type Storage interface {
	Import() ([]string, error)
	Export(rankings.Ranking) error
}
