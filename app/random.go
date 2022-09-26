package app

import "crypto/rand"

var (
	randomSecretLength      = 25
	randomSecretChars       = []byte(`0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`)
	randomSecretCharsLength = len(randomSecretChars)
)

func randomSecret() string {
	random := make([]byte, randomSecretLength)
	_, err := rand.Read(random)
	if err != nil {
		panic(err)
	}

	secret := make([]byte, randomSecretLength)
	for i := range random {
		secret[i] = randomSecretChars[random[i]%byte(randomSecretCharsLength)]
	}

	return string(secret)
}
