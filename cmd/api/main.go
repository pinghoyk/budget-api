package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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

	store := storage.New(db)

	// get
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "только GET", http.StatusMethodNotAllowed)
			return
		}

		users, err := store.GetAllUsers()
		if err != nil {
			log.Printf("Ошибка получения юзеров: %v", err)
			http.Error(w, "внутренняя ошибка", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	})

	// post
	http.HandleFunc("/add_user", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "только POST", http.StatusMethodNotAllowed)
			return
		}

		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			Name     string `json:"name"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "невалидный JSON", http.StatusBadRequest)
			return
		}

		if input.Email == "" || input.Password == "" || input.Name == "" {
			http.Error(w, "email, password, name обязательны", http.StatusBadRequest)
			return
		}

		err := store.AddUser(input.Email, input.Password, input.Name)
		if err != nil {
			log.Printf("Не добавил юзера: %v", err)
			http.Error(w, "не удалось создать юзера", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"status":"ok"}`))
		})

		log.Println("Сервер слушает :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
}

// пока все пишу тут, потом надо разделить файлы и кинут в internal