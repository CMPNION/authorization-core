package httpdelivery

import (
	"net/http"
)

func NewRouter(handler *AuthHandler, auth TokenAuthenticator) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/signup", handler.SignUp)
	mux.HandleFunc("/login", handler.SignIn)
	mux.Handle("/me", AuthMiddleware(auth)(http.HandlerFunc(handler.Me)))
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return logMiddleware(mux)
}
