package dp_golearn_to_rank

import (
	"testing"
)

func TestClient_FeatureStoreExists(t *testing.T) {
	c, err := NewClient()

	if err != nil {
		t.Fatal(err)
	}

	exists, err := c.FeatureService().DefaultFeatureStoreExists()

	if err != nil {
		t.Error(err)
	}

	if !exists {
		t.Error("feature store doesn't exist")
	}
}