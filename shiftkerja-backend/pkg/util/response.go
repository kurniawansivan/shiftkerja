package util

import (
	"encoding/json"
	"net/http"
	"shiftkerja-backend/internal/core/dto"
)

// RespondJSON sends a JSON response with the given status code
func RespondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// RespondSuccess sends a success response
func RespondSuccess(w http.ResponseWriter, message string, data interface{}) {
	response := dto.SuccessResponse{
		Message: message,
		Data:    data,
	}
	RespondJSON(w, http.StatusOK, response)
}

// RespondCreated sends a created response
func RespondCreated(w http.ResponseWriter, message string, data interface{}) {
	response := dto.SuccessResponse{
		Message: message,
		Data:    data,
	}
	RespondJSON(w, http.StatusCreated, response)
}

// RespondError sends an error response
func RespondError(w http.ResponseWriter, statusCode int, err string, message string) {
	response := dto.ErrorResponse{
		Error:   err,
		Message: message,
		Code:    statusCode,
	}
	RespondJSON(w, statusCode, response)
}

// RespondBadRequest sends a 400 Bad Request response
func RespondBadRequest(w http.ResponseWriter, message string) {
	RespondError(w, http.StatusBadRequest, "Bad Request", message)
}

// RespondUnauthorized sends a 401 Unauthorized response
func RespondUnauthorized(w http.ResponseWriter, message string) {
	RespondError(w, http.StatusUnauthorized, "Unauthorized", message)
}

// RespondForbidden sends a 403 Forbidden response
func RespondForbidden(w http.ResponseWriter, message string) {
	RespondError(w, http.StatusForbidden, "Forbidden", message)
}

// RespondNotFound sends a 404 Not Found response
func RespondNotFound(w http.ResponseWriter, message string) {
	RespondError(w, http.StatusNotFound, "Not Found", message)
}

// RespondInternalError sends a 500 Internal Server Error response
func RespondInternalError(w http.ResponseWriter, message string) {
	RespondError(w, http.StatusInternalServerError, "Internal Server Error", message)
}
