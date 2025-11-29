package dto

// CreateShiftRequest represents the request body for creating a shift
type CreateShiftRequest struct {
	Title       string  `json:"title" validate:"required,min=3,max=100"`
	Description string  `json:"description" validate:"max=500"`
	PayRate     float64 `json:"pay_rate" validate:"required,gt=0"`
	Lat         float64 `json:"lat" validate:"required,min=-90,max=90"`
	Lng         float64 `json:"lng" validate:"required,min=-180,max=180"`
}

// UpdateShiftRequest represents the request body for updating a shift
type UpdateShiftRequest struct {
	ID          int64   `json:"id" validate:"required"`
	Title       string  `json:"title" validate:"required,min=3,max=100"`
	Description string  `json:"description" validate:"max=500"`
	PayRate     float64 `json:"pay_rate" validate:"required,gt=0"`
	Lat         float64 `json:"lat" validate:"required,min=-90,max=90"`
	Lng         float64 `json:"lng" validate:"required,min=-180,max=180"`
	Status      string  `json:"status" validate:"required,oneof=OPEN FILLED"`
}

// ApplyShiftRequest represents the request body for applying to a shift
type ApplyShiftRequest struct {
	ShiftID int64 `json:"shift_id" validate:"required,gt=0"`
}

// UpdateApplicationStatusRequest represents the request body for updating application status
type UpdateApplicationStatusRequest struct {
	ApplicationID int64  `json:"application_id" validate:"required,gt=0"`
	Status        string `json:"status" validate:"required,oneof=ACCEPTED REJECTED"`
}

// ShiftResponse represents a shift in API responses
type ShiftResponse struct {
	ID          int64   `json:"id"`
	OwnerID     int64   `json:"owner_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	PayRate     float64 `json:"pay_rate"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
}

// ApplicationResponse represents an application in API responses
type ApplicationResponse struct {
	ID           int64   `json:"id"`
	ShiftID      int64   `json:"shift_id"`
	WorkerID     int64   `json:"worker_id"`
	Status       string  `json:"status"`
	CreatedAt    string  `json:"created_at"`
	ShiftTitle   string  `json:"shift_title,omitempty"`
	ShiftPayRate float64 `json:"shift_pay_rate,omitempty"`
	WorkerName   string  `json:"worker_name,omitempty"`
	WorkerEmail  string  `json:"worker_email,omitempty"`
}

// ErrorResponse represents an error in API responses
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}

// SuccessResponse represents a success message
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
