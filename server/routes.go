package server

import (
	"fmt"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/saurabh-anthwal/dummy/server/apis"
	"net/http"
	"time"
)

// Configure middlewares here.
var middlewares = []func(http.Handler) http.Handler{
	middleware.RealIP,
	middleware.Logger,
	middleware.Recoverer,
	middleware.Timeout(60 * time.Second), // timeout for request contexts.
}

// define your api routes here.
func createRoutes(router chi.Router) {
	// Set up our root handlers
	//router.Get("/", homePage)

	// Set up our API
	// router.Mount("/api/v1", v1.NewRouter())
	router.Route("/home", func(r chi.Router) {
		r.Get("/", apis.HelloWorld)
		r.Get("/name/{name}", apis.HelloName)
		r.Post("/user", apis.User)

	})
}

func configureRoutes() http.Handler {
	router := chi.NewRouter()

	// Set up our middlewares
	for _, mw := range middlewares {
		router.Use(mw)
	}

	createRoutes(router)
	return router
}

// maxBodyMiddleware sets a limit on request body.
func maxBodyMiddleware(sizeMB int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, 1024*sizeMB)
			err := r.ParseForm()
			if err != nil {
				zlog.Errorf("request body is invalid, error %v ", err)
				fmt.Fprintf(w, "request body is invalid, error: %v ", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// HelloWorld is a sample handler, that just prints hello world
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world from cripr 👋\n")
}
