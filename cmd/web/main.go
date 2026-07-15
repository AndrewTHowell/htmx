package main

import (
	"log/slog"
	"net/http"
	"os"

	"AndrewTHowell/htmx/assets"
)

// The application struct holds the dependencies needed for our handlers,
// including a htmlRenderer type.
type application struct {
	logger *slog.Logger
	html   *htmlRenderer
}

func securityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; img-src 'self' data:; style-src 'self' 'unsafe-inline'; object-src 'none'; base-uri 'self'; frame-ancestors 'none'")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		next.ServeHTTP(w, r)
	})
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Initialize a new htmlRenderer, parsing the base template and all partial
	// templates from assets/html into the shared template set.
	htmlRenderer, err := newHTMLRenderer(assets.HTMLFiles, "base.tmpl", "partials/*.tmpl")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Include the htmlRenderer in the application struct.
	app := &application{
		logger: logger,
		html:   htmlRenderer,
	}

	// Create a file server that serves the files from assets/static.
	fileserver := http.FileServerFS(assets.StaticFiles)

	// Register the application routes.
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static", fileserver))
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /gopher", app.gopher)
	mux.HandleFunc("GET /users", app.listUsers)
	mux.HandleFunc("GET /users/search", app.searchUsers)

	// Start the HTTP server.
	logger.Info("starting server", "port", 5051)
	err = http.ListenAndServe(":5051", securityMiddleware(mux))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
