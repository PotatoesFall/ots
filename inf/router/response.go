package router

import (
	"net/http"
)

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, `/`, http.StatusSeeOther)
}

func noContent(w http.ResponseWriter, code int) {
	w.Header().Set(`Content-Length`, `0`)
	w.WriteHeader(code)
}
