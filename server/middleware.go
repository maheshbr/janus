package server

import (
	"log"
	"net/http"
)

// corsHandler middleare to handle cors.
func corsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

// basicAuth middleware for handling basic auth request.
func basicAuth(username string, password string) handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			name, pass, ok := req.BasicAuth()
			if !ok || !(username == name && password == pass) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			h.ServeHTTP(w, req)
		})
	}
}

// recoverHandler middleware for panic recovery
func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func logHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method,r.Proto,r.URL)
		h.ServeHTTP(w, r)
	})
}