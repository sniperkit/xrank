package engine

import (
	"bytes"
	"fmt"
)

func checkUnitEqual(lhs RankUnit, rhs RankUnit) error {
	if lhs.ID != rhs.ID {
		return fmt.Errorf("ID not match, %d vs %d", lhs.ID, rhs.ID)
	}
	if lhs.Key != rhs.Key {
		return fmt.Errorf("Key not match, %d vs %d", lhs.Key, rhs.Key)
	}
	if bytes.Compare(lhs.Value, rhs.Value) != 0 {
		return fmt.Errorf("Value not match, %q vs %q", lhs.Value, rhs.Value)
	}
	return nil
}
