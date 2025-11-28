# ShiftKerja API Documentation

## Base URL
```
http://localhost:8080
```

## Authentication
All protected endpoints require JWT token in Authorization header:
```
Authorization: Bearer <token>
```

---

## Authentication Endpoints

### Register
**POST** `/register`

Create a new user account.

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "role": "worker" // or "business"
}
```

**Response:** `200 OK`
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "role": "worker"
  }
}
```

---

### Login
**POST** `/login`

Authenticate and receive JWT token.

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:** `200 OK`
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "role": "worker"
  }
}
```

---

## Shift Endpoints

### Get Nearby Shifts
**GET** `/shifts`
🔒 **Authentication Required**

Retrieve shifts within a specified radius from a location.

**Query Parameters:**
- `lat` (float, required): Latitude
- `lng` (float, required): Longitude  
- `rad` (float, optional): Radius in kilometers (default: 10)

**Example:**
```
GET /shifts?lat=-8.6478&lng=115.1385&rad=10
```

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "owner_id": 2,
    "title": "Barista at Canggu Coffee",
    "description": "Looking for experienced barista",
    "pay_rate": 75000,
    "lat": -8.6478,
    "lng": 115.1385,
    "status": "OPEN",
    "created_at": "2025-11-29T10:30:00Z"
  }
]
```

---

### Create Shift
**POST** `/shifts/create`
🔒 **Authentication Required** (Business only)

Post a new shift.

**Request Body:**
```json
{
  "title": "Barista at Canggu Coffee",
  "description": "Looking for experienced barista",
  "pay_rate": 75000,
  "lat": -8.6478,
  "lng": 115.1385
}
```

**Response:** `200 OK`
```json
{
  "shift": {
    "id": 1,
    "owner_id": 2,
    "title": "Barista at Canggu Coffee",
    "description": "Looking for experienced barista",
    "pay_rate": 75000,
    "lat": -8.6478,
    "lng": 115.1385,
    "status": "OPEN",
    "created_at": "2025-11-29T10:30:00Z"
  }
}
```

---

### Update Shift
**POST** `/shifts/update`
🔒 **Authentication Required** (Business owner only)

Update an existing shift.

**Request Body:**
```json
{
  "id": 1,
  "title": "Senior Barista at Canggu Coffee",
  "description": "Looking for experienced barista with 2+ years",
  "pay_rate": 85000,
  "lat": -8.6478,
  "lng": 115.1385,
  "status": "OPEN"
}
```

**Response:** `200 OK`
```json
{
  "status": "Shift updated successfully"
}
```

**Error Responses:**
- `403 Forbidden`: Not the shift owner
- `404 Not Found`: Shift doesn't exist
- `400 Bad Request`: Invalid data

---

### Delete Shift
**DELETE** `/shifts/delete`
🔒 **Authentication Required** (Business owner only)

Delete a shift and all its applications.

**Query Parameters:**
- `shift_id` (int, required): Shift ID to delete

**Example:**
```
DELETE /shifts/delete?shift_id=1
```

**Response:** `200 OK`
```json
{
  "status": "Shift deleted successfully"
}
```

**Error Responses:**
- `403 Forbidden`: Not the shift owner
- `404 Not Found`: Shift doesn't exist

---

### Get My Shifts
**GET** `/shifts/my-shifts`
🔒 **Authentication Required** (Business only)

Get all shifts posted by the authenticated business owner.

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "owner_id": 2,
    "title": "Barista at Canggu Coffee",
    "description": "Looking for experienced barista",
    "pay_rate": 75000,
    "lat": -8.6478,
    "lng": 115.1385,
    "status": "OPEN",
    "created_at": "2025-11-29T10:30:00Z"
  }
]
```

---

### Apply for Shift
**POST** `/shifts/apply`
🔒 **Authentication Required** (Worker only)

Apply for a shift.

**Request Body:**
```json
{
  "shift_id": 1
}
```

**Response:** `200 OK`
```json
{
  "status": "Applied successfully"
}
```

**Error Responses:**
- `403 Forbidden`: Not a worker
- `400 Bad Request`: Already applied or shift not available

---

## Application Endpoints

### Get My Applications
**GET** `/my-applications`
🔒 **Authentication Required** (Worker only)

Get all applications submitted by the authenticated worker.

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "shift_id": 1,
    "worker_id": 3,
    "status": "PENDING",
    "created_at": "2025-11-29T11:00:00Z",
    "shift_title": "Barista at Canggu Coffee",
    "shift_pay_rate": 75000
  }
]
```

---

### Get Shift Applications
**GET** `/shifts/applications`
🔒 **Authentication Required** (Business owner only)

Get all applications for a specific shift.

**Query Parameters:**
- `shift_id` (int, required): Shift ID

**Example:**
```
GET /shifts/applications?shift_id=1
```

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "shift_id": 1,
    "worker_id": 3,
    "worker_name": "John Doe",
    "worker_email": "john@example.com",
    "status": "PENDING",
    "created_at": "2025-11-29T11:00:00Z"
  }
]
```

---

### Update Application Status
**POST** `/shifts/applications/update`
🔒 **Authentication Required** (Business owner only)

Accept or reject an application.

**Request Body:**
```json
{
  "application_id": 1,
  "status": "ACCEPTED" // or "REJECTED"
}
```

**Response:** `200 OK`
```json
{
  "status": "Updated successfully"
}
```

**Note:** When an application is ACCEPTED, the shift status automatically changes to FILLED.

**Error Responses:**
- `403 Forbidden`: Not the shift owner
- `400 Bad Request`: Invalid status

---

### Delete Application (Withdraw)
**DELETE** `/my-applications/delete`
🔒 **Authentication Required** (Worker only)

Withdraw a pending application.

**Query Parameters:**
- `application_id` (int, required): Application ID to withdraw

**Example:**
```
DELETE /my-applications/delete?application_id=1
```

**Response:** `200 OK`
```json
{
  "status": "Application withdrawn successfully"
}
```

**Error Responses:**
- `403 Forbidden`: Not the application owner
- `400 Bad Request`: Can only withdraw PENDING applications

---

## WebSocket Endpoint

### Real-time Shift Updates
**WS** `/ws`

Connect to WebSocket for real-time shift updates.

**Connection:**
```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onmessage = (event) => {
  const shift = JSON.parse(event.data);
  console.log('New shift:', shift);
};
```

**Message Format:**
```json
{
  "id": "shift_001",
  "title": "Moving Barista",
  "lat": -8.6478,
  "lng": 115.1385,
  "pay_rate": 75000
}
```

---

## Health Check

### Health
**GET** `/health`

Check if the server is running.

**Response:** `200 OK`
```
ShiftKerja System Online
```

---

## Status Codes

- `200 OK`: Request successful
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

---

## Application Status Flow

```
PENDING → ACCEPTED (shift becomes FILLED)
       → REJECTED
```

## Shift Status

- `OPEN`: Accepting applications
- `FILLED`: Position filled, no longer accepting applications
- `CANCELLED`: Shift cancelled

---

## Complete CRUD Operations

### Shifts (Business Owners)
- ✅ **CREATE**: `POST /shifts/create`
- ✅ **READ**: `GET /shifts/my-shifts`
- ✅ **UPDATE**: `POST /shifts/update`
- ✅ **DELETE**: `DELETE /shifts/delete`

### Applications (Workers)
- ✅ **CREATE**: `POST /shifts/apply`
- ✅ **READ**: `GET /my-applications`
- ✅ **UPDATE**: N/A (Status updated by business owner)
- ✅ **DELETE**: `DELETE /my-applications/delete`

### Applications (Business Owners)
- ✅ **READ**: `GET /shifts/applications`
- ✅ **UPDATE**: `POST /shifts/applications/update`

---

## Example Usage Flow

### Worker Journey
1. Register: `POST /register` with role="worker"
2. Login: `POST /login`
3. Find shifts: `GET /shifts?lat=-8.6478&lng=115.1385&rad=10`
4. Apply: `POST /shifts/apply` with shift_id
5. Check status: `GET /my-applications`
6. Withdraw (if needed): `DELETE /my-applications/delete?application_id=1`

### Business Journey
1. Register: `POST /register` with role="business"
2. Login: `POST /login`
3. Create shift: `POST /shifts/create`
4. View shifts: `GET /shifts/my-shifts`
5. Check applications: `GET /shifts/applications?shift_id=1`
6. Accept/Reject: `POST /shifts/applications/update`
7. Edit shift: `POST /shifts/update`
8. Delete shift: `DELETE /shifts/delete?shift_id=1`

---

## Rate Limiting
Currently no rate limiting implemented. For production, consider adding rate limiting middleware.

## CORS
CORS is enabled for all origins in development. Update `cors_middleware.go` for production.
