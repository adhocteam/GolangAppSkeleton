package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	// This is inefficient - it reads the templates from the filesystem every
	// time. This makes it much easier to develop though, so we can edit our
	// templates and the changes will be reflected without having to restart
	// the app.
	t, err := template.ParseGlob("templates/*.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error %s", err.Error()), 500)
		return
	}

	err = t.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error %s", err.Error()), 500)
		return
	}
}

func logreq(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("path: %s", r.URL.Path)

		f(w, r)
	})
}

type App struct {
	Port       string
	StaticBase string
}

func (a App) Start() {
	if a.StaticBase == "/static" {
		log.Printf("serving static assets")
		http.Handle("/static/", logreq(staticHandler("static")))
	}
	http.Handle("/", logreq(a.index))
	addr := fmt.Sprintf(":%s", a.Port)
	log.Printf("Starting app on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (a App) index(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", struct {
		Name       string
		StaticBase string
	}{
		Name:       "Sonic The Hedgehog",
		StaticBase: a.StaticBase,
	})
}

func staticHandler(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/static/", http.FileServer(http.Dir(dir))).ServeHTTP(w, r)
	}
}

func env(key, adefault string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return adefault
	}
	return val
}

func main() {
	server := App{
		Port:       env("PORT", "8080"),
		StaticBase: env("STATIC_BASE", "/static"),
	}
	server.Start()
}
