package store

import (
	"testing"
)

func TestHexToObjectID(t *testing.T) {
	hexObjID := "599e8545da144128eb510159"
	fakeObjID := "qwessdf"

	bsonObjID, _ := hexToObjectID(hexObjID)
	if bsonObjID == "" {
		t.Fatal("Expected to have bson object ID")
	}

	_, err := hexToObjectID(fakeObjID)
	if err == nil {
		t.Errorf("Expected to get error: invalid metadata ID, got: %v", err)
	}
}
