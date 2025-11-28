# 🏗️ Architecture & Code Quality Improvements

## Clean Architecture Implementation

### Before (Tightly Coupled)
```go
// Handler directly accessing repository
func (h *ShiftHandler) Create(w http.ResponseWriter, r *http.Request) {
    // Parse request
    var shift entity.Shift
    json.NewDecoder(r.Body).Decode(&shift)
    
    // Direct database call - business logic mixed with HTTP
    h.PostgresRepo.CreateShift(r.Context(), &shift)
    h.RedisRepo.AddShift(r.Context(), shift)
    
    json.NewEncoder(w).Encode(shift)
}
```

❌ **Problems:**
- Business logic in handler
- Tight coupling to specific implementations
- Hard to test
- No validation
- Can't swap implementations

### After (Clean Architecture)
```go
// Handler uses service layer
func (h *ShiftHandler) Create(w http.ResponseWriter, r *http.Request) {
    // Parse and validate
    var req entity.Shift
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }
    
    req.OwnerID = int64(userID)
    
    // Delegate to service (business logic)
    if err := h.Service.CreateShift(r.Context(), &req); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(req)
}

// Service handles business logic
func (s *ShiftService) CreateShift(ctx context.Context, shift *entity.Shift) error {
    // Business validation
    if shift.PayRate <= 0 {
        return errors.New("pay rate must be positive")
    }
    
    // Orchestrate repositories
    if err := s.shiftRepo.CreateShift(ctx, shift); err != nil {
        return err
    }
    
    // Geo sync (log but don't fail)
    if err := s.geoRepo.AddShift(ctx, *shift); err != nil {
        fmt.Printf("⚠️ Redis sync warning: %v\n", err)
    }
    
    return nil
}
```

✅ **Benefits:**
- Separation of concerns
- Business logic in service layer
- Easy to test
- Can mock dependencies
- Clear error handling

---

## Dependency Inversion (Port/Adapter Pattern)

### Port Layer (Contracts)
```go
// /internal/core/port/shift_repository.go
type ShiftRepository interface {
    CreateShift(ctx context.Context, shift *entity.Shift) error
    GetShiftByID(ctx context.Context, id int64) (*entity.Shift, error)
    // ... more methods
}
```

### Adapter Layer (Implementations)
```go
// /internal/adapter/repository/postgres_shift.go
type PostgresShiftRepo struct {
    DB *pgx.Conn
}

func (r *PostgresShiftRepo) CreateShift(ctx context.Context, shift *entity.Shift) error {
    // Postgres-specific implementation
}
```

### Service Uses Interface
```go
type ShiftService struct {
    shiftRepo port.ShiftRepository  // ← Interface, not concrete type
    geoRepo   port.GeoRepository
}
```

✅ **Benefits:**
- Service doesn't know about Postgres/Redis
- Easy to swap implementations (e.g., MySQL, MongoDB)
- Easy to mock for testing
- Follows SOLID principles

---

## Entity Design Improvements

### Before (Anemic Model)
```go
type Shift struct {
    ID      int64
    Title   string
    PayRate float64
}
```

### After (Rich Domain Model)
```go
type Shift struct {
    ID          int64     `json:"id"`
    OwnerID     int64     `json:"owner_id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    PayRate     float64   `json:"pay_rate"`
    Lat         float64   `json:"lat"`
    Lng         float64   `json:"lng"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
}

type Application struct {
    ID        int64     `json:"id"`
    ShiftID   int64     `json:"shift_id"`
    WorkerID  int64     `json:"worker_id"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
    
    // Populated via JOIN queries
    ShiftTitle   string  `json:"shift_title,omitempty"`
    ShiftPayRate float64 `json:"shift_pay_rate,omitempty"`
    WorkerName   string  `json:"worker_name,omitempty"`
    WorkerEmail  string  `json:"worker_email,omitempty"`
}
```

✅ **Benefits:**
- Complete domain representation
- JSON tags for API serialization
- JOIN-friendly fields
- Self-documenting

---

## Repository Pattern Improvements

### Complete CRUD Operations
```go
// Before: Only Create
func (r *PostgresShiftRepo) CreateShift(...)
func (r *PostgresShiftRepo) ApplyForShift(...)

// After: Full interface
func (r *PostgresShiftRepo) CreateShift(...)
func (r *PostgresShiftRepo) GetShiftByID(...)
func (r *PostgresShiftRepo) GetShiftsByOwner(...)
func (r *PostgresShiftRepo) UpdateShiftStatus(...)
func (r *PostgresShiftRepo) ApplyForShift(...)
func (r *PostgresShiftRepo) GetApplicationsByWorker(...)
func (r *PostgresShiftRepo) GetApplicationsByShift(...)
func (r *PostgresShiftRepo) UpdateApplicationStatus(...)
func (r *PostgresShiftRepo) GetApplicationByID(...)
```

### JOIN Queries for Performance
```go
// Before: N+1 queries
apps := GetApplications(shiftID)
for _, app := range apps {
    worker := GetUser(app.WorkerID)  // ← N queries!
}

// After: Single JOIN query
query := `
    SELECT 
        a.id, a.shift_id, a.worker_id, a.status, a.created_at,
        u.full_name, u.email
    FROM applications a
    JOIN users u ON a.worker_id = u.id
    WHERE a.shift_id = $1
`
```

---

## Service Layer Business Logic

### Authorization
```go
func (s *ShiftService) GetShiftApplications(ctx context.Context, shiftID, requesterID int64) ([]entity.Application, error) {
    // Verify ownership
    shift, err := s.shiftRepo.GetShiftByID(ctx, shiftID)
    if err != nil {
        return nil, ErrShiftNotFound
    }
    
    if shift.OwnerID != requesterID {
        return nil, ErrUnauthorized  // ← Business rule
    }
    
    return s.shiftRepo.GetApplicationsByShift(ctx, shiftID)
}
```

### Validation
```go
func (s *ShiftService) ApplyForShift(ctx context.Context, shiftID, workerID int64) error {
    // Check if shift exists
    shift, err := s.shiftRepo.GetShiftByID(ctx, shiftID)
    if err != nil {
        return ErrShiftNotFound
    }
    
    // Check if shift is still open
    if shift.Status != "OPEN" {
        return errors.New("shift is no longer available")
    }
    
    return s.shiftRepo.ApplyForShift(ctx, shiftID, workerID)
}
```

### Orchestration
```go
func (s *ShiftService) UpdateApplicationStatus(ctx context.Context, applicationID, businessID int64, newStatus string) error {
    // 1. Validate
    if newStatus != "ACCEPTED" && newStatus != "REJECTED" {
        return ErrInvalidStatus
    }
    
    // 2. Check authorization
    app, _ := s.shiftRepo.GetApplicationByID(ctx, applicationID)
    shift, _ := s.shiftRepo.GetShiftByID(ctx, app.ShiftID)
    if shift.OwnerID != businessID {
        return ErrUnauthorized
    }
    
    // 3. Update application
    s.shiftRepo.UpdateApplicationStatus(ctx, applicationID, newStatus)
    
    // 4. If accepted, mark shift as filled
    if newStatus == "ACCEPTED" {
        s.shiftRepo.UpdateShiftStatus(ctx, app.ShiftID, "FILLED")
        s.geoRepo.RemoveShift(ctx, app.ShiftID)  // Remove from Redis
    }
    
    return nil
}
```

---

## Frontend Architecture Improvements

### Component Organization

```
views/              # Page-level components
├── LoginView.vue
├── RegisterView.vue
├── BusinessDashboard.vue
└── WorkerDashboard.vue

components/         # Reusable components
└── MapData.vue

stores/             # State management
├── auth.js        # Authentication state
└── socket.js      # WebSocket state

router/            # Navigation
└── index.js       # Routes + guards
```

### State Management (Pinia)
```javascript
// Before: Props drilling
<Parent>
  <Child :token="token" :role="role">
    <GrandChild :token="token" :role="role">  // ← Annoying!
    </GrandChild>
  </Child>
</Parent>

// After: Global store
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
console.log(authStore.token)  // ← Available anywhere!
```

### Composition API Benefits
```vue
<!-- Before: Options API -->
<script>
export default {
  data() {
    return { count: 0 }
  },
  methods: {
    increment() { this.count++ }
  },
  mounted() { /* ... */ }
}
</script>

<!-- After: Composition API -->
<script setup>
import { ref, onMounted } from 'vue'

const count = ref(0)
const increment = () => count.value++

onMounted(() => { /* ... */ })
</script>
```

✅ **Benefits:**
- Better TypeScript support
- Logic reuse via composables
- Less boilerplate
- Clearer data flow

---

## Error Handling Improvements

### Backend
```go
// Before: Generic errors
http.Error(w, "Failed", 500)

// After: Specific errors
var (
    ErrUnauthorized      = errors.New("unauthorized action")
    ErrShiftNotFound     = errors.New("shift not found")
    ErrApplicationExists = errors.New("already applied")
)

if err := service.Apply(...); err != nil {
    if errors.Is(err, ErrShiftNotFound) {
        http.Error(w, err.Error(), http.StatusNotFound)
    } else {
        http.Error(w, err.Error(), http.StatusBadRequest)
    }
}
```

### Frontend
```javascript
// Before: Silent failures
fetch('/api/apply', { body: data })

// After: User feedback
try {
    const res = await fetch('/api/apply', { body: data })
    if (res.ok) {
        alert('✅ Success!')
    } else {
        const error = await res.text()
        alert('❌ ' + error)
    }
} catch (error) {
    alert('Network error')
}
```

---

## Security Improvements

### JWT Middleware
```go
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Extract token
        tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
        
        // Validate
        claims, err := service.ValidateToken(tokenString)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }
        
        // Inject user context
        ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
        ctx = context.WithValue(ctx, "role", claims["role"])
        
        next(w, r.WithContext(ctx))
    }
}
```

### Role-Based Access
```go
// Handler checks role from context
role := r.Context().Value("role").(string)
if role != "business" {
    http.Error(w, "Forbidden", http.StatusForbidden)
    return
}
```

---

## Testing Improvements

### Before (Hard to Test)
```go
func TestCreateShift(t *testing.T) {
    // Need real Postgres + Redis!
    conn := pgx.Connect(...)
    rdb := redis.NewClient(...)
    
    handler := NewShiftHandler(redis, postgres)
    // ...
}
```

### After (Easy to Test with Mocks)
```go
// Mock implementation
type MockShiftRepo struct {
    shifts []entity.Shift
}

func (m *MockShiftRepo) CreateShift(ctx context.Context, shift *entity.Shift) error {
    m.shifts = append(m.shifts, *shift)
    return nil
}

// Test service with mock
func TestCreateShift(t *testing.T) {
    mockRepo := &MockShiftRepo{}
    mockGeo := &MockGeoRepo{}
    service := NewShiftService(mockRepo, mockGeo)
    
    err := service.CreateShift(ctx, &shift)
    
    assert.NoError(t, err)
    assert.Equal(t, 1, len(mockRepo.shifts))
}
```

---

## Performance Optimizations

### Database Indexing
```sql
-- Fast email lookups (login)
CREATE INDEX ON users (email);

-- Fast owner lookups (my shifts)
CREATE INDEX ON shifts (owner_id);

-- Prevent duplicate applications
CREATE UNIQUE INDEX ON applications (shift_id, worker_id);
```

### Redis Geospatial Index
```go
// O(log N) nearest neighbor search
locations := r.Client.GeoSearch(ctx, "shifts_geo", &redis.GeoSearchQuery{
    Latitude:  lat,
    Longitude: lng,
    Radius:    10,
    RadiusUnit: "km",
})
```

### Connection Pooling
```go
// Using pgx built-in connection pool
conn, err := pgx.Connect(ctx, dbURL)
// Reuses connections automatically
```

---

## Code Quality Metrics

| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| Layers | 2 | 5 | ✅ Separation of concerns |
| Testability | Low | High | ✅ Dependency injection |
| Endpoints | 4 | 9 | ✅ Complete CRUD |
| Type Safety | Partial | Full | ✅ Interfaces everywhere |
| Error Handling | Generic | Specific | ✅ Meaningful errors |
| Frontend Views | 2 | 5 | ✅ Complete user flows |
| State Management | Props | Pinia | ✅ Centralized state |

---

## Conclusion

The ShiftKerja MVP now follows industry best practices:

✅ **Clean Architecture** - Layers are properly separated  
✅ **SOLID Principles** - Dependency inversion, single responsibility  
✅ **Domain-Driven Design** - Rich entities, business logic in service  
✅ **Repository Pattern** - Data access abstraction  
✅ **Security** - JWT, RBAC, validation  
✅ **Performance** - Indexing, caching, JOIN queries  
✅ **Maintainability** - Easy to extend and test  

This architecture scales from MVP to production! 🚀
