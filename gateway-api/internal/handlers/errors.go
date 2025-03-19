package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	oz "github.com/go-ozzo/ozzo-validation/v4"
)

func WriteError(w http.ResponseWriter, err string, code int) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err})
}

func WriteValidationError(w http.ResponseWriter, err error) {
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
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
