package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"shiftkerja-backend/internal/adapter/handler"
	"shiftkerja-backend/internal/adapter/repository"
	"shiftkerja-backend/internal/core/service"

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

	// --- 4. SERVICES (Business Logic Layer) ---
	shiftService := service.NewShiftService(pgShiftRepo, redisRepo)

	// --- 5. HANDLERS & ROUTES ---

	// WebSocket Hub (created first to pass to handlers)
	wsHub := handler.NewHub()

	// A. Shift Handlers (with WebSocket hub for broadcasting)
	shiftHandler := handler.NewShiftHandler(shiftService, wsHub)

	// Shift Routes
	http.HandleFunc("/shifts", handler.AuthMiddleware(shiftHandler.GetNearby))
	http.HandleFunc("/shifts/create", handler.AuthMiddleware(shiftHandler.Create))
	http.HandleFunc("/shifts/update", handler.AuthMiddleware(shiftHandler.UpdateShift))
	http.HandleFunc("/shifts/delete", handler.AuthMiddleware(shiftHandler.DeleteShift))
	http.HandleFunc("/shifts/apply", handler.AuthMiddleware(shiftHandler.Apply))
	http.HandleFunc("/shifts/my-shifts", handler.AuthMiddleware(shiftHandler.GetMyShifts))
	http.HandleFunc("/shifts/applications", handler.AuthMiddleware(shiftHandler.GetShiftApplications))
	http.HandleFunc("/shifts/applications/update", handler.AuthMiddleware(shiftHandler.UpdateApplicationStatus))
	
	// Worker Routes
	http.HandleFunc("/my-applications", handler.AuthMiddleware(shiftHandler.GetMyApplications))
	http.HandleFunc("/my-applications/delete", handler.AuthMiddleware(shiftHandler.DeleteApplication))

	// B. Auth Handlers
	authHandler := handler.NewAuthHandler(userRepo)
	http.HandleFunc("/register", authHandler.Register)
	http.HandleFunc("/login", authHandler.Login)

	// C. WebSocket Endpoint
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