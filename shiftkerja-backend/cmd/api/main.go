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

	// --- 3. SEED DATA (Temporary) ---
	repo := repository.NewRedisGeoRepo(rdb)

	mockShift := entity.Shift{
		ID:      "shift_001",
		Title:   "Barista at Canggu Coffee",
		Lat:     -8.6478,
		Lng:     115.1385,
		PayRate: 75000,
	}

	// We seed silently here to ensure data exists for the map
	_ = repo.AddShift(context.Background(), mockShift)

	// --- 4. HTTP Server Setup ---
	
	// A. Initialize the Handler
	// We inject the 'repo' so the handler can access the database
	shiftHandler := handler.NewShiftHandler(repo)

	// B. Register the Routes
	// When a user hits "/shifts", run the GetNearby function
	http.HandleFunc("/shifts", shiftHandler.GetNearby) 

	// C. Health Check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ShiftKerja System Online"))
	})

	// D. Start Server
	fmt.Println("🚀 ShiftKerja Backend starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error:", err)
	}
}