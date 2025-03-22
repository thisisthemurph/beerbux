package send

import (
	"encoding/json"
	"errors"
	"net/http"

	oz "github.com/go-ozzo/ozzo-validation/v4"
)

// JSON sends the given JSON payload with the given status code.
//
// If there is an error sending the payload, a generic Error is returned.
func JSON(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		Error(w, "Failed to respond", http.StatusInternalServerError)
	}
}

// Error sends a simple error JSON payload with the given status code.
func Error(w http.ResponseWriter, err string, code int) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err})
}

// ValidationError sends a JSON payload detailing ozzo-validation validation errors.
// The JSON payload contains an errors map detailing the fields and their errors.
//
// This method sets the status code to 400 by default.
//
// If there is an error sending the ValidationError payload, a generic Error is returned.
func ValidationError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	if err == nil {
		return
	}

	result := map[string]interface{}{
		"errors": make(map[string]string),
	}

	var validationErrors oz.Errors
	if errors.As(err, &validationErrors) {
		fieldErrs := result["errors"].(map[string]string)
		for field, e := range validationErrors {
			fieldErrs[field] = e.Error()
		}
	} else {
		// Handle any unexpected errors
		result["error"] = err.Error()
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		Error(w, "Failed to encode validation error", http.StatusInternalServerError)
	}
}
