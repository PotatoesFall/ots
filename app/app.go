package app

import (
	"errors"
	"sync"

	"github.com/PotatoesFall/ots/domain"
)

// Errors
var (
	ErrSecretNotFound = errors.New(`secret not found`)
)

// App is able to perform all application logic. It should always be instantiated using app.New(...)
type App interface {
	// NewSecret creates a new secret and returns the ID
	NewSecret(domain.NewSecret) int

	// NewRandomSecret generates a random secret, creates it and returns the ID
	NewRandomSecret(message string) int

	// ClaimSecret claims a secret.
	ClaimSecret(id int) (domain.Secret, error)

	// SecretExists checks whether a secret exists (and is not claimed)
	SecretExists(id int) bool
}

type app struct {
	sync.Mutex // locking for claims

	repo Repo
}

func (a *app) NewSecret(s domain.NewSecret) int {
	return a.repo.Create(s)
}

func (a *app) NewRandomSecret(message string) int {
	return a.repo.Create(domain.NewSecret{
		Message: message,
		Content: randomSecret(),
	})
}

func (a *app) ClaimSecret(id int) (domain.Secret, error) {
	a.Lock()
	defer a.Unlock()

	secret, err := a.repo.Claim(id)
	if err != nil && !errors.Is(err, ErrSecretNotFound) {
		panic(err)
	}

	return secret, err
}

func (a *app) SecretExists(id int) bool {
	return a.repo.Exists(id)
}

type Repo interface {
	Exists(id int) bool
	Create(domain.NewSecret) int
	Claim(id int) (domain.Secret, error)
}

func New(repo Repo) App {
	return &app{repo: repo}
}
