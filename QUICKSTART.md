# 🚀 ShiftKerja - Quick Start Guide

## Prerequisites
- Docker Desktop installed and running
- Go 1.22+ installed
- Node.js 18+ installed
- Make (optional, but helpful)

## 🏃‍♂️ Quick Start (5 Minutes)

### Step 1: Start Infrastructure
```bash
cd /Users/admin/Documents/Ivan/projects/shiftkerja
docker compose up -d
```
✅ This starts PostgreSQL (with PostGIS) and Redis

### Step 2: Apply Database Migrations
```bash
make migrateup
```
Or without Make:
```bash
migrate -path db/migration -database "postgres://postgres:password123@localhost:5432/shiftkerja?sslmode=disable" up
```

### Step 3: Start Backend
```bash
cd shiftkerja-backend
go run cmd/api/main.go
```
✅ Backend runs on `http://localhost:8080`

### Step 4: Start Frontend (New Terminal)
```bash
cd shiftkerja-frontend
npm install  # First time only
npm run dev
```
✅ Frontend runs on `http://localhost:5173`

---

## 🎭 Demo Scenario

### Create Business Account
1. Open `http://localhost:5173`
2. Click **"Register"**
3. Fill form:
   - Name: `Cafe Owner`
   - Email: `business@test.com`
   - Password: `password123`
   - Role: **Business**
4. Auto-login → Map view

### Post a Shift
1. Click **"📊 Dashboard"** button (top right)
2. Click **"➕ Post New Shift"**
3. Fill form:
   - Title: `Barista Needed`
   - Description: `Morning shift at our cafe`
   - Pay Rate: `80000`
   - Lat: `-8.6478`
   - Lng: `115.1385`
4. Click **"Create Shift"**
5. ✅ Shift appears in your dashboard

### Apply as Worker (Incognito Window)
1. Open new **Incognito/Private window**
2. Go to `http://localhost:5173`
3. Click **"Register"**
4. Fill form:
   - Name: `John Worker`
   - Email: `worker@test.com`
   - Password: `password123`
   - Role: **Worker**
5. Auto-login → Map shows the shift pin
6. Click the **pin** on map
7. Click **"📋 Apply Now"**
8. ✅ Application submitted!

### View Application (Worker)
1. Click **"📊 Dashboard"**
2. See your application with status: **PENDING**

### Accept Application (Business Window)
1. Switch back to business window
2. Refresh dashboard (or click Dashboard again)
3. See application from `John Worker`
4. Click **"✓ Accept"** button
5. ✅ Shift status changes to **FILLED**
6. ✅ Application status changes to **ACCEPTED**

### Verify (Worker Window)
1. Switch to worker window
2. Refresh dashboard
3. ✅ Status shows **ACCEPTED** with success message

---

## 🧪 Test WebSocket (Real-time Feature)

1. Open **two browser windows** side-by-side (both logged in)
2. In either window, click **"📍 Test WS"** button
3. ✅ New marker appears **instantly** on both maps
4. This demonstrates real-time location broadcasting

---

## 📡 API Testing (cURL)

### Register Business
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "api-business@test.com",
    "password": "test123",
    "full_name": "API Business",
    "role": "business"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "api-business@test.com",
    "password": "test123"
  }'
```
Copy the `token` from response.

### Create Shift
```bash
curl -X POST http://localhost:8080/shifts/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "title": "Waiter Needed",
    "description": "Evening shift",
    "pay_rate": 75000,
    "lat": -8.65,
    "lng": 115.14
  }'
```

### Get Nearby Shifts
```bash
curl -X GET "http://localhost:8080/shifts?lat=-8.65&lng=115.14&rad=10" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

## 🛑 Stop Everything

```bash
# Stop backend: Ctrl+C in backend terminal
# Stop frontend: Ctrl+C in frontend terminal

# Stop databases
docker compose down
```

---

## 🐛 Troubleshooting

### Backend won't start
```bash
# Check if Go modules are installed
cd shiftkerja-backend
go mod tidy
go mod download
```

### Frontend won't start
```bash
cd shiftkerja-frontend
rm -rf node_modules
npm install
```

### Database connection error
```bash
# Check if Docker is running
docker ps

# Restart containers
docker compose down
docker compose up -d

# Wait 5 seconds, then retry backend
```

### Port already in use
```bash
# Backend (8080)
lsof -ti:8080 | xargs kill -9

# Frontend (5173)
lsof -ti:5173 | xargs kill -9
```

### Migration fails
```bash
# Check if database exists
psql postgres://postgres:password123@localhost:5432/shiftkerja -c "SELECT 1"

# If fails, create database
docker exec -it shiftkerja-postgres-1 psql -U postgres -c "CREATE DATABASE shiftkerja;"
```

---

## 📂 Project Structure

```
shiftkerja/
├── docker-compose.yml          # Infrastructure setup
├── Makefile                    # Useful commands
├── db/
│   └── migration/              # SQL migrations
├── shiftkerja-backend/         # Go API
│   ├── cmd/api/main.go         # Entry point
│   ├── internal/
│   │   ├── core/               # Business logic
│   │   │   ├── entity/         # Domain models
│   │   │   ├── port/           # Interfaces
│   │   │   └── service/        # Business services
│   │   └── adapter/            # Infrastructure
│   │       ├── handler/        # HTTP handlers
│   │       └── repository/     # Database access
└── shiftkerja-frontend/        # Vue 3 SPA
    ├── src/
    │   ├── views/              # Page components
    │   ├── components/         # Reusable components
    │   ├── stores/             # Pinia state
    │   └── router/             # Vue Router
```

---

## ✅ Checklist

- [ ] Docker running
- [ ] Postgres & Redis containers up (`docker ps`)
- [ ] Migrations applied (`make migrateup`)
- [ ] Backend running on `:8080`
- [ ] Frontend running on `:5173`
- [ ] Can register business
- [ ] Can register worker
- [ ] Business can post shift
- [ ] Worker can see shift on map
- [ ] Worker can apply
- [ ] Business can accept/reject
- [ ] WebSocket test works

---

## 🎉 Success!

If all checkboxes are ✅, your ShiftKerja MVP is fully operational!

**Next Steps:**
- Read `MVP_COMPLETION.md` for detailed feature documentation
- Review `README.md` for architecture overview
- Check backend logs for request/response details
- Explore the code to understand Clean Architecture implementation

Happy coding! 🚀
