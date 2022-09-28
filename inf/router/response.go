package router

import (
	"net/http"
)

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, `/`, http.StatusSeeOther)
}

func writeString(w http.ResponseWriter, code int, body string) {
	w.Header().Set(`Content-Type`, `text/plain`)
	w.WriteHeader(code)
	_, err := w.Write([]byte(body))
	if err != nil {
		panic(err)
	}
}

func noContent(w http.ResponseWriter, code int) {
	w.Header().Set(`Content-Length`, `0`)
	w.WriteHeader(code)
}
