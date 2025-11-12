package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"errors"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/pinghoyk/budget-api/internal/storage"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден!")
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/budget.db"
		log.Printf("DB_PATH не задан — используем значение по умолчанию: %s", dbPath)
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

	log.Printf("База данных инициализирована: %s", dbPath)

	store := storage.New(db)

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes
	r.Get("/users", getUsersHandler(store))
	r.Post("/users", addUserHandler(store))
	r.Delete("/users/{id}", deleteUserHandler(store))
	r.Get("/users/{id}", getUserHandler(store))

	// Запуск сервера
	log.Println("Сервер слушает :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// === HANDLERS ===

func getUsersHandler(store *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := store.GetAllUsers()
		if err != nil {
			log.Printf("Ошибка получения юзеров: %v", err)
			http.Error(w, "внутренняя ошибка", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(users); err != nil {
			log.Printf("Ошибка сериализации JSON: %v", err)
			http.Error(w, "ошибка при кодировании", http.StatusInternalServerError)
			return
		}
	}
}

func addUserHandler(store *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			Name     string `json:"name"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "невалидный JSON", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if input.Email == "" || input.Password == "" || input.Name == "" {
			http.Error(w, "email, password, name обязательны", http.StatusBadRequest)
			return
		}

		err := store.AddUser(input.Email, input.Password, input.Name)
		if err != nil {
			log.Printf("Не удалось добавить юзера: %v", err)
			http.Error(w, "не удалось создать юзера", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}

func deleteUserHandler(store *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := chi.URLParam(r, "id")
		if userIDStr == "" {
			http.Error(w, "ID пользователя обязателен", http.StatusBadRequest)
			return
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			http.Error(w, "ID должен быть целым числом", http.StatusBadRequest)
			return
		}

		if userID <= 0 {
			http.Error(w, "ID должен быть положительным", http.StatusBadRequest)
			return
		}

		err = store.DeleteUser(userID)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				http.Error(w, "пользователь не найден", http.StatusNotFound)
				return
			}
			log.Printf("Ошибка удаления пользователя (ID=%d): %v", userID, err)
			http.Error(w, "внутренняя ошибка", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}

func getUserHandler(store *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := chi.URLParam(r, "id")
		if userIDStr == "" {
			http.Error(w, "ID пользователя обязателен", http.StatusBadRequest)
			return
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			http.Error(w, "ID должен быть целым числом", http.StatusBadRequest)
			return
		}

		if userID <= 0 {
			http.Error(w, "ID должен быть положительным", http.StatusBadRequest)
			return
		}

		user, err := store.GetUserByID(userID)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				http.Error(w, "пользователь не найден", http.StatusNotFound)
				return
			}
			log.Printf("Ошибка получения пользователя (ID=%d): %v", userID, err)
			http.Error(w, "внутренняя ошибка", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(user); err != nil {
			log.Printf("Ошибка сериализации JSON для пользователя (ID=%d): %v", userID, err)
			http.Error(w, "ошибка при кодировании", http.StatusInternalServerError)
			return
		}
	}
}