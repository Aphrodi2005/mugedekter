// handlers.go
package main

import (
	"Cinema/pkg/models"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func (app *application) showAllMovies(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/all-movies" {
		http.NotFound(w, r)
		return
	}

	movies, err := app.movies.Latest(10)
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
	// Получение данных о фильме и месте из формы запроса
	movieID := r.FormValue("movie_id")
	seat := r.FormValue("seat")

	// Преобразование строки в ObjectID для использования в модели
	objectID, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Вызов метода модели для покупки билета
	err = app.movies.BuyTicket(objectID, seat)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Перенаправление пользователя на страницу с фильмами
	http.Redirect(w, r, "/all-movies", http.StatusSeeOther)
}

func (app *application) returnTicket(w http.ResponseWriter, r *http.Request) {
	// Получение данных о фильме и месте из формы запроса
	movieID := r.FormValue("movie_id")
	seat := r.FormValue("seat")

	// Преобразование строки в ObjectID для использования в модели
	objectID, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Вызов метода модели для возврата билета
	err = app.movies.ReturnTicket(objectID, seat)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Перенаправление пользователя на страницу с фильмами
	http.Redirect(w, r, "/all-movies", http.StatusSeeOther)
}

func (app *application) retrieveMovie(w http.ResponseWriter, r *http.Request) {
	// Получение данных о фильме из формы запроса
	movieID := r.FormValue("movie_id")

	// Преобразование строки в ObjectID для использования в модели
	objectID, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Вызов метода модели для получения информации о фильме
	movie, err := app.movies.Get(objectID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Отображение информации о фильме
	err = app.render(w, r, "retrieve_movie.page.tmpl", &templateData{
		Movie: movie,
	})
	if err != nil {
		app.serverError(w, err)
	}
}
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Home page logic
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Initialize a slice containing the paths to the two files. Note that the
	// home.page.tmpl file must be the *first* file in the slice.
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	// Using template.ParseFiles() function to read the template file
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// Using the Execute() method on the template set to write the template
	// content as the response body.
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
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
	ratingStr := r.PostForm.Get("rating")

	// Преобразование строки в целое число
	rating, err := strconv.Atoi(ratingStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.movies.Create(title, genre, rating)
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

	id, err := primitive.ObjectIDFromHex(r.PostForm.Get("id"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	title := r.PostForm.Get("title")
	genre := r.PostForm.Get("genre")
	// Добавьте другие поля, которые могут быть у фильма

	err = app.movies.Update(title, genre, 0, id) // Передайте оценку (rating) как 0, если у вас нет этого поля
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

	id := r.FormValue("id")
	if id == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.movies.Delete(objectID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/all-movies", http.StatusSeeOther)
}
