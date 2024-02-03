// main.go

package main

import (
	"context"
	"fmt"
	"github.com/Aphrodi2005/database" // Замените на фактический путь к вашему файлу database.go
)

func main() {
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

	// Ваш код для создания сервера
	fmt.Println("Сервер успешно создан.")
	// ...
}
