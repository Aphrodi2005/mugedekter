// handlers.go
package main

import (
	"Cinema/pkg/models"
	"errors"
	"net/http"
	"strconv"
)

func (app *application) showAllMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := app.movies.All()
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.render(w, r, "all_movies.page.tmpl", &templateData{
		Movies: movies,
	})
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) buyTicket(w http.ResponseWriter, r *http.Request) {
	// Логика для покупки билета
}

func (app *application) returnTicket(w http.ResponseWriter, r *http.Request) {
	// Логика для возврата билета
}

func (app *application) retrieveMovie(w http.ResponseWriter, r *http.Request) {
	// Логика для получения информации о фильме
}
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Логика для возврата билета
}
func (app *application) contacts(w http.ResponseWriter, r *http.Request) {
	// Логика для возврата билета
}

func (app *application) addMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}

	title := r.PostForm.Get("title")
	genre := r.PostForm.Get("genre")
	// Добавьте другие поля, которые могут быть у фильма

	err = app.movies.Add(title, genre) // Используйте метод вашей модели для добавления фильма
	if errors.Is(err, models.ErrDuplicate) {
		app.clientError(w, http.StatusBadRequest)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/all-movies", http.StatusSeeOther)
}

func (app *application) updateMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}

	id, err := strconv.Atoi(r.PostForm.Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	title := r.PostForm.Get("title")
	genre := r.PostForm.Get("genre")
	// Добавьте другие поля, которые могут быть у фильма

	err = app.movies.Update(title, genre, id) // Используйте метод вашей модели для обновления фильма

	if errors.Is(err, models.ErrDuplicate) {
		app.clientError(w, http.StatusBadRequest)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/all-movies", http.StatusSeeOther)
}

func (app *application) deleteMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.movies.Delete(id) // Используйте метод вашей модели для удаления фильма
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/all-movies", http.StatusSeeOther)
}
