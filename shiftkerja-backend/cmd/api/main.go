package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"shiftkerja-backend/internal/core/entity"
	"shiftkerja-backend/internal/adapter/repository"

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

	// --- 2.5. SEED DATA (The Missing Link) ---
	// We initialize the repo using the connection we just made
	repo := repository.NewRedisGeoRepo(rdb)

	// Create a fake shift in Bali (Canggu area)
	mockShift := entity.Shift{
		ID:      "shift_001",
		Title:   "Barista at Canggu Coffee",
		Lat:     -8.6478,
		Lng:     115.1385,
		PayRate: 75000,
	}

	fmt.Println("🌱 Seeding Mock Shift to Redis...")
	// We use a background context here for simplicity in the main function
	if err := repo.AddShift(context.Background(), mockShift); err != nil {
		fmt.Printf("❌ Failed to seed: %v\n", err)
	} else {
		fmt.Println("✅ Mock Shift Seeded! (ID: shift_001)")
	}

	// Optional: Quick self-test to prove it's there
	found, _ := repo.FindNearby(context.Background(), -8.64, 115.13, 10)
	fmt.Printf("🔍 Self-Test: Found %d shifts nearby!\n", len(found))

	// --- 3. Start HTTP Server ---
	fmt.Println("🚀 ShiftKerja Backend starting on port 8080...")
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ShiftKerja System Online"))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error:", err)
	}
}