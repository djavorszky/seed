package example

import (
	"log"
	"net/http"
)

type MyServer struct{}

func (s *MyServer) LoggerMw(next http.Handler) http.Handler {
	// Anything you add here will be executed once, during startup.
	// The returned http.Handler will be able to access these variables
	// thanks to closure.

	prefix := "[NAME] - "
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(prefix, r.RemoteAddr, r.Method, r.RequestURI)

		next.ServeHTTP(w, r)
	})
}

func (s *MyServer) Index() http.HandlerFunc {
	// Anything you add here will be executed once, during startup.
	// The returned http.HandlerFunc will be able to access these variables
	// thanks to closure.

	defaultMsg := []byte("I'm alive!")
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		_, err := w.Write(defaultMsg)
		if err != nil {
			panic(err)
		}
	}
}
