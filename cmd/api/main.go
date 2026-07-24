package main

import (
	"fmt"
	"net/http"
	"time"

	"RATE_LIMITER_MIDDLEWARE/internal/limiter"
)

func main() {
	cfg := limiter.Config{
		Rate:       1 * time.Second,
		MaxTokens:  3,
		CleanupTTL: 10 * time.Second,
	}

	store := limiter.NewMemoryStore(cfg)
	defer store.Close()

	rateLimiter := limiter.NewMiddleware(store)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "200 OK: Доступ разрешен")
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      rateLimiter.Handler(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Println("Сервер запущен на http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Ошибка сервера: %v\n", err)
	}
}
