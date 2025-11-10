package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/pinghoyk/budget-api/internal/auth"
	"github.com/pinghoyk/budget-api/internal/handlers"
	"github.com/pinghoyk/budget-api/internal/storage"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден!")
	}

	// Validate required environment variables
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET не установлен в переменных окружения")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN не установлен в переменных окружения")
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

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(store)
	accountHandler := handlers.NewAccountHandler(store)
	categoryHandler := handlers.NewCategoryHandler(store)
	transactionHandler := handlers.NewTransactionHandler(store)
	budgetHandler := handlers.NewBudgetHandler(store)
	statsHandler := handlers.NewStatsHandler(store)

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// Public routes
	r.Route("/api", func(r chi.Router) {
		// Health check
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		// Authentication routes (public)
		r.Post("/auth/telegram", authHandler.TelegramLogin)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(auth.Middleware(jwtSecret))

			// Current user
			r.Get("/auth/me", authHandler.GetCurrentUser)

			// Accounts
			r.Route("/accounts", func(r chi.Router) {
				r.Get("/", accountHandler.GetAccounts)
				r.Post("/", accountHandler.CreateAccount)
				r.Get("/{id}", accountHandler.GetAccount)
				r.Put("/{id}", accountHandler.UpdateAccount)
				r.Delete("/{id}", accountHandler.DeleteAccount)
			})

			// Categories
			r.Route("/categories", func(r chi.Router) {
				r.Get("/", categoryHandler.GetCategories)
				r.Post("/", categoryHandler.CreateCategory)
				r.Get("/{id}", categoryHandler.GetCategory)
				r.Put("/{id}", categoryHandler.UpdateCategory)
				r.Delete("/{id}", categoryHandler.DeleteCategory)
			})

			// Transactions
			r.Route("/transactions", func(r chi.Router) {
				r.Get("/", transactionHandler.GetTransactions)
				r.Post("/", transactionHandler.CreateTransaction)
				r.Get("/{id}", transactionHandler.GetTransaction)
				r.Put("/{id}", transactionHandler.UpdateTransaction)
				r.Delete("/{id}", transactionHandler.DeleteTransaction)
			})

			// Budgets
			r.Route("/budgets", func(r chi.Router) {
				r.Get("/", budgetHandler.GetBudgets)
				r.Post("/", budgetHandler.CreateBudget)
				r.Get("/{id}", budgetHandler.GetBudget)
				r.Get("/{id}/status", budgetHandler.GetBudgetStatus)
				r.Put("/{id}", budgetHandler.UpdateBudget)
				r.Delete("/{id}", budgetHandler.DeleteBudget)
			})

			// Statistics
			r.Route("/stats", func(r chi.Router) {
				r.Get("/category-summary", statsHandler.GetCategorySummary)
				r.Get("/monthly-balance", statsHandler.GetMonthlyBalance)
			})
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Сервер слушает :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}