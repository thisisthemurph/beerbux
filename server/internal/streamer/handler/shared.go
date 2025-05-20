package handler

import "net/http"

func getFromQuery(w http.ResponseWriter, r *http.Request, key string) (string, bool) {
	value := r.URL.Query().Get(key)
	if value == "" {
		http.Error(w, key+" is required", http.StatusBadRequest)
		return "", false
	}
	return value, true
}

func getUserID(w http.ResponseWriter, r *http.Request) (string, bool) {
	return getFromQuery(w, r, "user_id")
}

func getSessionID(w http.ResponseWriter, r *http.Request) (string, bool) {
	return getFromQuery(w, r, "session_id")
}

func setServerSentEventHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
}
