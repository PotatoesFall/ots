package app

import "golang.org/x/crypto/bcrypt"

const cost = 12

func hash(in string) string {
	out, err := bcrypt.GenerateFromPassword([]byte(in), cost)
	if err != nil {
		panic(err)
	}

	return string(out)
}
