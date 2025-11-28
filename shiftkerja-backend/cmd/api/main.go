package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

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
	// We use the address from docker-compose
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", 
		DB:       0,  
	})

	// PING Redis to verify connection
	fmt.Println("⏳ Connecting to Redis...")
	redisCtx, redisCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer redisCancel()

	if err := rdb.Ping(redisCtx).Err(); err != nil {
		fmt.Printf("❌ Unable to connect to Redis: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Connected to Redis successfully!")

	// --- 3. Start HTTP Server ---
	fmt.Println("🚀 ShiftKerja Backend starting on port 8080...")
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ShiftKerja System Online"))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error:", err)
	}
}