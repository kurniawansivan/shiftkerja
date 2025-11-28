# 🎯 ShiftKerja - Command Reference

## Quick Commands

### Start Everything
```bash
# 1. Start infrastructure
docker compose up -d

# 2. Apply migrations
make migrateup

# 3. Start backend (terminal 1)
cd shiftkerja-backend && go run cmd/api/main.go

# 4. Start frontend (terminal 2)
cd shiftkerja-frontend && npm run dev
```

### Stop Everything
```bash
# Stop backend: Ctrl+C in terminal 1
# Stop frontend: Ctrl+C in terminal 2
docker compose down
```

---

## Docker Commands

```bash
# Start containers
docker compose up -d

# Stop containers
docker compose down

# View logs
docker compose logs -f

# Restart containers
docker compose restart

# Remove everything (including volumes)
docker compose down -v
```

---

## Database Commands

### Migrations
```bash
# Apply all migrations
make migrateup

# Rollback one migration
make migratedown

# Create new migration
migrate create -ext sql -dir db/migration -seq migration_name

# Manual migration (without Make)
migrate -path db/migration \
  -database "postgres://postgres:password123@localhost:5432/shiftkerja?sslmode=disable" \
  up
```

### Direct Database Access
```bash
# Connect to PostgreSQL
docker exec -it shiftkerja-postgres-1 psql -U postgres -d shiftkerja

# Common queries
\dt                    # List tables
\d users              # Describe users table
SELECT * FROM users;  # Query users
\q                    # Quit
```

### Redis Access
```bash
# Connect to Redis
docker exec -it shiftkerja-redis-1 redis-cli

# Common commands
KEYS *                        # List all keys
GEORADIUS shifts_geo ...      # Geo search
GET shift:101                 # Get shift data
FLUSHALL                      # Clear everything (careful!)
```

---

## Backend Commands

### Development
```bash
cd shiftkerja-backend

# Run server
go run cmd/api/main.go

# Install dependencies
go mod tidy
go mod download

# Format code
go fmt ./...

# Run tests (when implemented)
go test ./...

# Build binary
go build -o bin/api cmd/api/main.go

# Run binary
./bin/api
```

### Debugging
```bash
# Enable verbose logging
GODEBUG=gctrace=1 go run cmd/api/main.go

# Check module status
go mod verify

# List dependencies
go list -m all
```

---

## Frontend Commands

### Development
```bash
cd shiftkerja-frontend

# Install dependencies
npm install

# Start dev server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Lint code
npm run lint

# Format code
npm run format
```

### Package Management
```bash
# Add dependency
npm install package-name

# Remove dependency
npm uninstall package-name

# Update dependencies
npm update

# Check outdated packages
npm outdated
```

---

## API Testing with cURL

### Authentication
```bash
# Register Business
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "business@test.com",
    "password": "password123",
    "full_name": "Business Owner",
    "role": "business"
  }'

# Register Worker
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "worker@test.com",
    "password": "password123",
    "full_name": "Worker Name",
    "role": "worker"
  }'

# Login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "business@test.com",
    "password": "password123"
  }'

# Save token
export TOKEN="eyJhbGc..."
```

### Shifts
```bash
# Create Shift (Business)
curl -X POST http://localhost:8080/shifts/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "Barista Needed",
    "description": "Morning shift",
    "pay_rate": 80000,
    "lat": -8.6478,
    "lng": 115.1385
  }'

# Get Nearby Shifts
curl -X GET "http://localhost:8080/shifts?lat=-8.6478&lng=115.1385&rad=10" \
  -H "Authorization: Bearer $TOKEN"

# Get My Shifts (Business)
curl -X GET http://localhost:8080/shifts/my-shifts \
  -H "Authorization: Bearer $TOKEN"
```

### Applications
```bash
# Apply for Shift (Worker)
curl -X POST http://localhost:8080/shifts/apply \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"shift_id": 1}'

# Get My Applications (Worker)
curl -X GET http://localhost:8080/my-applications \
  -H "Authorization: Bearer $TOKEN"

# Get Shift Applications (Business)
curl -X GET "http://localhost:8080/shifts/applications?shift_id=1" \
  -H "Authorization: Bearer $TOKEN"

# Update Application Status (Business)
curl -X POST http://localhost:8080/shifts/applications/update \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "application_id": 1,
    "status": "ACCEPTED"
  }'
```

### Health Check
```bash
curl http://localhost:8080/health
```

---

## Troubleshooting Commands

### Port Issues
```bash
# Check what's using port 8080
lsof -ti:8080

# Kill process on port 8080
lsof -ti:8080 | xargs kill -9

# Check what's using port 5173
lsof -ti:5173

# Kill process on port 5173
lsof -ti:5173 | xargs kill -9
```

### Database Issues
```bash
# Check if Postgres is running
docker ps | grep postgres

# View Postgres logs
docker logs shiftkerja-postgres-1

# Restart Postgres
docker restart shiftkerja-postgres-1

# Connect and verify database
docker exec -it shiftkerja-postgres-1 psql -U postgres -c "SELECT 1"
```

### Redis Issues
```bash
# Check if Redis is running
docker ps | grep redis

# View Redis logs
docker logs shiftkerja-redis-1

# Test Redis connection
docker exec -it shiftkerja-redis-1 redis-cli PING
```

### Backend Issues
```bash
# Check Go version
go version

# Verify module integrity
cd shiftkerja-backend && go mod verify

# Re-download dependencies
go mod download

# Clean cache
go clean -modcache
```

### Frontend Issues
```bash
# Clear node_modules and reinstall
cd shiftkerja-frontend
rm -rf node_modules package-lock.json
npm install

# Clear build cache
npm run build -- --force

# Check Node version
node --version
npm --version
```

---

## Useful Aliases (Add to ~/.zshrc)

```bash
# ShiftKerja shortcuts
alias sk-start='docker compose up -d'
alias sk-stop='docker compose down'
alias sk-backend='cd shiftkerja-backend && go run cmd/api/main.go'
alias sk-frontend='cd shiftkerja-frontend && npm run dev'
alias sk-logs='docker compose logs -f'
alias sk-db='docker exec -it shiftkerja-postgres-1 psql -U postgres -d shiftkerja'
alias sk-redis='docker exec -it shiftkerja-redis-1 redis-cli'
alias sk-migrate='make migrateup'
```

After adding, reload:
```bash
source ~/.zshrc
```

---

## Environment Variables

### Backend
```bash
# Create .env file (future enhancement)
DATABASE_URL=postgres://postgres:password123@localhost:5432/shiftkerja
REDIS_URL=localhost:6379
JWT_SECRET=SUPER_SECRET_KEY_DO_NOT_SHARE
PORT=8080
```

### Frontend
```bash
# Create .env file
VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080/ws
```

---

## Git Commands (Version Control)

```bash
# Initialize (if not done)
git init
git add .
git commit -m "Complete MVP implementation"

# Create feature branch
git checkout -b feature/payment-integration

# Push to remote
git remote add origin https://github.com/yourusername/shiftkerja.git
git push -u origin main

# Useful checks
git status
git log --oneline
git diff
```

---

## Performance Monitoring

### Backend
```bash
# CPU & Memory profiling
go tool pprof http://localhost:8080/debug/pprof/profile
go tool pprof http://localhost:8080/debug/pprof/heap

# Check goroutines
curl http://localhost:8080/debug/pprof/goroutine?debug=1
```

### Database
```bash
# Check slow queries (Postgres)
docker exec -it shiftkerja-postgres-1 psql -U postgres -d shiftkerja -c \
  "SELECT * FROM pg_stat_statements ORDER BY total_time DESC LIMIT 10;"

# Check connection count
docker exec -it shiftkerja-postgres-1 psql -U postgres -d shiftkerja -c \
  "SELECT count(*) FROM pg_stat_activity;"
```

### Redis
```bash
# Monitor commands
docker exec -it shiftkerja-redis-1 redis-cli MONITOR

# Get stats
docker exec -it shiftkerja-redis-1 redis-cli INFO

# Memory usage
docker exec -it shiftkerja-redis-1 redis-cli INFO memory
```

---

## Quick Reference URLs

- **Frontend:** http://localhost:5173
- **Backend API:** http://localhost:8080
- **Health Check:** http://localhost:8080/health
- **WebSocket:** ws://localhost:8080/ws
- **PostgreSQL:** localhost:5432
- **Redis:** localhost:6379

---

## One-Line Setup

```bash
docker compose up -d && sleep 5 && make migrateup && (cd shiftkerja-backend && go run cmd/api/main.go) & (cd shiftkerja-frontend && npm run dev)
```

Press `Ctrl+C` twice to stop both servers.
