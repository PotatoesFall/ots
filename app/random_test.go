package app

import (
	"testing"

	"golang.org/x/exp/slices"
)

func TestRandomSecret(t *testing.T) {
	secret := randomSecret()

	if len(secret) != randomSecretLength {
		t.Errorf("%s has incorrect length %d, should be %d", secret, len(secret), randomSecretCharsLength)
	}

	for _, b := range []byte(secret) {
		if slices.Index(randomSecretChars, b) == -1 {
			t.Errorf(`%s contains illegal character %c`, secret, b)
		}
	}
}
