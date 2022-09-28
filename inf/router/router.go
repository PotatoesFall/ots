package router

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/PotatoesFall/ots/app"
	"github.com/PotatoesFall/ots/inf/router/form"
	"github.com/PotatoesFall/ots/inf/router/views"
	"github.com/go-chi/chi"
)

func Make(a app.App) http.Handler {
	r := chi.NewRouter()
	r.MethodNotAllowed(methodNotAllowedHandler)
	r.Use(recoverMiddleware)

	r.Get(`/`, getIndex)
	r.Get(`/secret/{id}`, getSecret(a))
	r.Get(`/random`, getRandom(a))
	r.Post(`/new`, newSecret(a))
	r.Post(`/claim/{id}`, claimSecret(a))

	return r
}

func readForm[T any](w http.ResponseWriter, r *http.Request) (T, bool) {
	var t T

	if err := r.ParseForm(); err != nil {
		noContent(w, http.StatusBadRequest)
		return t, false
	}

	if err := form.Read(r.Form, &t); err != nil {
		noContent(w, http.StatusBadRequest)
		return t, false
	}

	return t, true
}

func getIndex(w http.ResponseWriter, req *http.Request) {
	views.Render(w, `index`, nil)
}

func getSecret(a app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, `id`)
		id, err := strconv.Atoi(idStr)
		hash, exists := a.CheckSecret(id)
		if err != nil || !exists {
			views.Render(w, `notfound`, nil)
			return
		}

		views.Render(w, `secret`, views.GetSecretResponse{ID: id, Hash: hash})
	}
}

func getRandom(a app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		secret := a.GenerateRandomSecret()
		writeString(w, http.StatusOK, secret)
	}
}

func newSecret(a app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, ok := readForm[NewSecretRequest](w, r)
		if !ok {
			return
		}

		id := a.NewSecret(req.Secret())

		http.Redirect(w, r, `/secret/`+strconv.Itoa(id), http.StatusFound)
	}
}

func claimSecret(a app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, `id`)
		id, err := strconv.Atoi(idStr)
		if err != nil {
			noContent(w, http.StatusBadRequest)
			return
		}

		secret, err := a.ClaimSecret(id)
		if errors.Is(err, app.ErrSecretNotFound) {
			views.Render(w, `notfound`, nil)
			return
		}
		if err != nil {
			panic(err)
		}

		views.Render(w, `claim`, views.NewClaimSecretResponse(secret))
	}
}
