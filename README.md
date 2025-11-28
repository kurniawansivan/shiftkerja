# 🚀 ShiftKerja: Real-Time Geo-Spatial Shift Marketplace

ShiftKerja is a high-concurrency, location-based platform connecting businesses with shift workers in real-time. Think of it as **"Uber for Part-time Jobs."**

It utilizes a **Dual-Write Architecture** where persistent data lives in **PostgreSQL** and live geospatial data is indexed in **Redis** for blazing-fast "Nearest Neighbor" searches. Real-time updates are pushed via **WebSockets** to a **Vue 3** frontend.

## 🏗️ Tech Stack

### Backend (The Core)

  * **Language:** Golang 1.22+
  * **Framework:** Standard Library + Gorilla Mux/WebSocket
  * **Architecture:** Clean Architecture (Onion Layers)
  * **Security:** JWT (JSON Web Tokens) + BCrypt + Middleware
  * **Drivers:** `pgx/v5` (Postgres), `go-redis/v9` (Redis)

### Frontend (The Interface)

  * **Framework:** Vue 3 (Composition API)
  * **Styling:** TailwindCSS v4
  * **State Management:** Pinia (Auth & Socket Stores)
  * **Maps:** Leaflet.js + OpenStreetMap

### Infrastructure (The Foundation)

  * **Database:** PostgreSQL 18 + PostGIS Extension
  * **Cache/Geo Engine:** Redis 7
  * **Containerization:** Docker & Docker Compose
  * **Migrations:** `golang-migrate`

-----

## 🏛️ System Architecture

The project follows a strict **Clean Architecture** to ensure separation of concerns:

```text
/cmd/api           # Application Entry Point (Main, Config, Router)
/internal
  /core            # Pure Business Logic (No SQL/HTTP imports)
    /entity        # Domain Models (User, Shift)
    /service       # Complex Business Rules (Token Generation)
  /adapter         # Infrastructure Layer
    /handler       # HTTP Handlers & WebSocket Hub
    /repository    # Database Implementations (Postgres & Redis)
/db/migration      # SQL Migration Files
```

### The Real-Time Geo Engine

1.  **Ingestion:** When a Business posts a shift, it is saved to **Postgres** (Source of Truth) AND **Redis** (Geospatial Index).
2.  **Search:** Workers query **Redis** (`GEORADIUS`) to find jobs within 10km in milliseconds.
3.  **Live Stream:** Workers transmit GPS coordinates via **WebSockets**. The Go server broadcasts these updates to active clients for live map tracking.

-----

## 💾 Database Schema

### 1\. Users Table (`users`)

Stores authentication and role data.

  * `id`: BIGSERIAL (PK)
  * `email`: VARCHAR (Unique, Indexed)
  * `password_hash`: VARCHAR (Bcrypt)
  * `role`: VARCHAR ('worker' | 'business' | 'admin')

### 2\. Shifts Table (`shifts`)

Stores the job postings.

  * `id`: BIGSERIAL (PK)
  * `owner_id`: BIGINT (FK -\> users.id)
  * `lat` / `lng`: FLOAT8 (Synced to Redis)
  * `status`: VARCHAR ('OPEN', 'FILLED')

### 3\. Applications Table (`applications`)

Links Workers to Shifts.

  * `id`: BIGSERIAL (PK)
  * `shift_id`: BIGINT (FK -\> shifts.id)
  * `worker_id`: BIGINT (FK -\> users.id)
  * `status`: VARCHAR ('PENDING')
  * **Unique Constraint:** `(shift_id, worker_id)` prevents double applying.

-----

## 🔌 API Documentation

### 🔐 Authentication

| Method | Endpoint | Auth? | Description | Payload |
| :--- | :--- | :--- | :--- | :--- |
| **POST** | `/register` | No | Create new account | `{email, password, role, full_name}` |
| **POST** | `/login` | No | Get JWT Token | `{email, password}` |

### 🌍 Geospatial & Shifts

| Method | Endpoint | Auth? | Description | Payload |
| :--- | :--- | :--- | :--- | :--- |
| **GET** | `/shifts` | **Yes** | Find nearby shifts | Query Params: `?lat=-8.6&lng=115.1&rad=10` |
| **POST** | `/shifts/create` | **Yes** | Post a new shift | `{title, pay_rate, lat, lng, description}` |
| **POST** | `/shifts/apply` | **Yes** | Apply for a job | `{shift_id}` |

### ⚡ Real-Time (WebSocket)

  * **URL:** `ws://localhost:8080/ws`
  * **Protocol:** JSON
  * **Function:** Connects client to the Broadcast Hub.
  * **Incoming Message:** `{"lat": -8.6, "lng": 115.1, "status": "moving"}`
  * **Outgoing Message:** Server broadcasts received payload to all connected clients.

-----

## 🚀 Getting Started

### Prerequisites

  * Docker Desktop
  * Go 1.22+
  * Node.js 18+

### 1\. Start Infrastructure

Spin up Postgres (PostGIS) and Redis.

```bash
docker compose up -d
```

### 2\. Apply Database Migrations

Create the tables in Postgres.

```bash
make migrateup
# OR if you don't have Make:
# migrate -path db/migration -database "postgres://postgres:password123@localhost:5432/shiftkerja?sslmode=disable" up
```

### 3\. Start Backend

Run the Go API server on port 8080.

```bash
cd shiftkerja-backend
go run cmd/api/main.go
```

### 4\. Start Frontend

Run the Vue 3 development server.

```bash
cd shiftkerja-frontend
npm run dev
```

-----

## 🧪 Testing Flow

1.  **Register a Business:**
      * `POST /register` with `role: "business"`.
2.  **Login:**
      * `POST /login` -\> Copy the `token`.
3.  **Post a Shift:**
      * `POST /shifts/create` with the token.
      * Check logs: "Saved to Postgres" & "Synced to Redis".
4.  **Register a Worker:**
      * `POST /register` with `role: "worker"`.
5.  **View Map:**
      * Worker logs in on Frontend.
      * Map loads -\> Fetches `/shifts` -\> Shows the pin created in step 3.
6.  **Real-Time Movement:**
      * Open two browser windows.
      * Click "Move Me" in one window.
      * Watch the pin move instantly in the other window via WebSockets.

-----

## 🔮 Future Roadmap (Sprint 5+)

  * [ ] **Payment Gateway:** Integrate a dummy payment gateaway
  * [ ] **Notification System:** Push notifications when a worker applies.
  * [ ] **Admin Dashboard:** React-based admin panel for moderation.
  * [ ] **CI/CD:** GitHub Actions pipeline for automated testing and deployment.