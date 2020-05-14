package confusables

import (
	"testing"
)

func TestNormalizeAlice(t *testing.T) {
	username1 := "Alice" // Uses all ASCII characters, like you would naively expect
	username2 := "Αlice" // Uses a Greek letter A
	if username1 == username2 {
		t.Errorf("Usernames are equal before normalization. This should never happen.")
	}

	username2 = Normalize(username1)
	if username1 != username2 {
		t.Errorf("Normalization did not make the usernames equal.")
	}
}
