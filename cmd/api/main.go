package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/pinghoyk/budget-api/internal/storage"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден!")
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "../budget-data/budget.db"
		log.Printf("DB_PATH не задан - используем значение по умолчанию: %s", dbPath)
	}

	dataDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Не удалось создать директорию %s: %v", dataDir, err)
	}

	db, err := storage.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Не удалось инициализировать БД: %v", err)
	}
	defer db.Close()

	log.Printf("База данных создана :%s", dbPath)
	log.Println("Запустили сервер")

}

// пока все пишу тут, потом надо разделить файлы и кинут в internal