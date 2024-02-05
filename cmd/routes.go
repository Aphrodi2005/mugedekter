// routes.go
package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := mux.NewRouter()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/buy-ticket", app.buyTicket).Methods("POST")
	mux.HandleFunc("/return-ticket", app.returnTicket).Methods("POST")

	mux.HandleFunc("/all-movies", app.showAllMovies) // Добавлено

	mux.HandleFunc("/add-movie", app.addMovie).Methods("POST")
	mux.HandleFunc("/update-movie", app.updateMovie).Methods("POST")
	mux.HandleFunc("/delete-movie", app.deleteMovie).Methods("DELETE")
	mux.HandleFunc("/retrieve-movie", app.retrieveMovie)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))

	return mux
}
