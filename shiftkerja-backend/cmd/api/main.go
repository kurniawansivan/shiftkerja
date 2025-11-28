package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"shiftkerja-backend/internal/adapter/handler"
	"shiftkerja-backend/internal/adapter/repository"
	"shiftkerja-backend/internal/core/entity"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

func main() {
	// --- 1. Database Setup (Postgres) ---
	dbURL := "postgres://postgres:password123@localhost:5432/shiftkerja"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("⏳ Connecting to Postgres...")
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		fmt.Printf("❌ Unable to connect to Postgres: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	if err := conn.Ping(ctx); err != nil {
		fmt.Printf("❌ Postgres Ping failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Connected to Postgres successfully!")

	// --- 2. Cache Setup (Redis) ---
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	fmt.Println("⏳ Connecting to Redis...")
	redisCtx, redisCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer redisCancel()

	if err := rdb.Ping(redisCtx).Err(); err != nil {
		fmt.Printf("❌ Unable to connect to Redis: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Connected to Redis successfully!")

	// --- 3. REPOSITORIES ---
	redisRepo := repository.NewRedisGeoRepo(rdb)
	pgShiftRepo := repository.NewPostgresShiftRepo(conn)
	userRepo := repository.NewPostgresUserRepo(conn)

	// --- 4. SEED DATA ---
	mockShift := entity.Shift{
		ID:      101, // int64
		Title:   "Barista at Canggu Coffee",
		Lat:     -8.6478,
		Lng:     115.1385,
		PayRate: 75000,
	}
	_ = redisRepo.AddShift(context.Background(), mockShift)

	// --- 5. HANDLERS & ROUTES ---

	// A. Shift Handlers
	shiftHandler := handler.NewShiftHandler(redisRepo, pgShiftRepo)

	// Public/Protected Shift Routes
	// 1. Get Nearby (Map)
	http.HandleFunc("/shifts", handler.AuthMiddleware(shiftHandler.GetNearby))
	// 2. Create Shift (Business)
	http.HandleFunc("/shifts/create", handler.AuthMiddleware(shiftHandler.Create))
	// 3. Apply for Shift (Worker) - 👈 NEW ENDPOINT ADDED HERE
	http.HandleFunc("/shifts/apply", handler.AuthMiddleware(shiftHandler.Apply))

	// B. Auth Handlers
	authHandler := handler.NewAuthHandler(userRepo)
	http.HandleFunc("/register", authHandler.Register)
	http.HandleFunc("/login", authHandler.Login)

	// C. WebSocket Setup
	wsHub := handler.NewHub()
	http.HandleFunc("/ws", wsHub.HandleWS)

	// D. Health Check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ShiftKerja System Online"))
	})

	// --- 6. START SERVER ---
	fmt.Println("🚀 ShiftKerja Backend starting on port 8080...")

	// Wrap the default router with CORS Middleware
	router := handler.CORSMiddleware(http.DefaultServeMux)

	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println("Error:", err)
	}
}