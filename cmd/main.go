// main.go
package main

import (
	"Cinema/pkg"
	"Cinema/pkg/models"
	"context"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	movies        *models.MovieModel
	templateCache map[string]*template.Template
}

func main() {

	addr := flag.String("addr", ":8080", "HTTP network address") // Изменил порт на 8080, но вы можете использовать ваш
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Подключение к базе данных
	err := database.Connect()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := database.Client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Инициализация необходимых компонентов приложения
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(), // Подставьте сюда свои обработчики
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
