package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/PotatoesFall/ots/app"
	"github.com/PotatoesFall/ots/domain"
)

type Repo struct {
	sync.Mutex // locking for claims

	db *sql.DB
}

func New(db *sql.DB) *Repo {
	migrate(db)
	return &Repo{db: db}
}

func (r *Repo) Exists(id int) bool {
	row := r.db.QueryRow(`SELECT 1 FROM secret s WHERE s.id = $1;`, id)
	err := row.Scan(new(int))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		panic(err)
	}

	return err == nil
}

func (r *Repo) Create(s domain.NewSecret) int {
	var id int
	err := r.db.QueryRow(`INSERT INTO secret (message, content, hash) VALUES ($1, $2, $3) RETURNING id;`, s.Message, s.Content, s.Hash).Scan(&id)
	if err != nil {
		panic(err)
	}

	return id
}

func (r *Repo) ClaimHash(id int) (string, error) {
	r.Lock()
	defer r.Unlock()

	tx, err := r.db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback() //nolint:errcheck

	row := tx.QueryRow(`SELECT s.hash FROM secret s WHERE s.ID = $1;`, id)
	var hash string
	if err := row.Scan(&hash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", app.ErrSecretNotFound
		}

		panic(err)
	}

	res, err := tx.Exec(`UPDATE secret SET hash = '' WHERE ID = $1;`, id)
	if err != nil {
		panic(err)
	}
	if affected, err := res.RowsAffected(); affected != 1 || err != nil {
		panic(fmt.Errorf(`%d rows affected: %w`, affected, err))
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}

	return hash, nil
}

func (r *Repo) Claim(id int) (domain.Secret, error) {
	r.Lock()
	defer r.Unlock()

	tx, err := r.db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback() //nolint:errcheck

	row := tx.QueryRow(`SELECT s.message, s.content FROM secret s WHERE s.ID = $1;`, id)
	secret := domain.Secret{ID: id}
	if err := row.Scan(&secret.Message, &secret.Content); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return secret, app.ErrSecretNotFound
		}

		panic(err)
	}

	res, err := tx.Exec(`DELETE FROM secret s WHERE s.ID = $1;`, id)
	if err != nil {
		panic(err)
	}
	if affected, err := res.RowsAffected(); affected != 1 || err != nil {
		panic(fmt.Errorf(`%d rows affected: %w`, affected, err))
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}

	return secret, nil
}
