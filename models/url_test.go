package models

import "testing"

func TestURL_GenerateKey(t *testing.T) {
	u := &URL{}
	lastKey := ""
	t.Run("key should not be empty", func(t *testing.T) {
		u.GenerateKey()
		lastKey = u.Key

		if u.Key == "" {
			t.Errorf("key is empty")
		}
	})

	t.Run("key should be changed", func(t *testing.T) {
		u.GenerateKey()
		if u.Key == lastKey {
			t.Errorf("key is not changed")
		}
	})
}
