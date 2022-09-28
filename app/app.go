package app

import (
	"errors"

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

	// GenerateRandomSecret makes a random secret
	GenerateRandomSecret() string

	// ClaimSecret claims a secret.
	ClaimSecret(id int) (domain.Secret, error)

	// CheckSecret checks whether a secret exists (and is not claimed). It retuns the hash if not claimed yet.
	CheckSecret(id int) (hash string, exists bool)
}

type app struct {
	repo Repo
}

func (a *app) NewSecret(s domain.NewSecret) int {
	if s.Content == `` {
		s.Content = randomSecret()
		s.Hash = hash(s.Content)
	}

	return a.repo.Create(s)
}

func (a *app) GenerateRandomSecret() string {
	return randomSecret()
}

func (a *app) ClaimSecret(id int) (domain.Secret, error) {
	secret, err := a.repo.Claim(id)
	if err != nil && !errors.Is(err, ErrSecretNotFound) {
		panic(err)
	}

	return secret, err
}

func (a *app) CheckSecret(id int) (string, bool) {
	hash, err := a.repo.ClaimHash(id)
	if err != nil && !errors.Is(err, ErrSecretNotFound) {
		panic(err)
	}

	return hash, err == nil
}

type Repo interface {
	Exists(id int) bool
	Create(domain.NewSecret) int
	Claim(id int) (domain.Secret, error)
	ClaimHash(id int) (string, error)
}

func New(repo Repo) App {
	return &app{repo: repo}
}
