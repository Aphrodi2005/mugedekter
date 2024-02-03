// routes.go
package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := mux.NewRouter()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/buy-ticket", app.buyTicket)
	mux.HandleFunc("/return-ticket", app.returnTicket)
	mux.HandleFunc("/contacts", app.contacts)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))

	return mux
}
