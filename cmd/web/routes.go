package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	//mux.Handle("/",http.HandlerFunc(home))
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir(".\\ui\\static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	//mux.HandleFunc("/download/file.zip", app.downloadHandler)

	return mux
}