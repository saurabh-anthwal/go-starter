package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func jsonResponse(w http.ResponseWriter, data interface{}, c int) {
	dj, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		zlog.Errorw("Error creating JSON response", "error", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c)
	fmt.Fprintf(w, "%s", dj)
}

func errJsonResponse(w http.ResponseWriter, err error) {
	zlog.Error(err)
	resp := struct {
		Error string
	}{err.Error()}
	jsonResponse(w, resp, http.StatusInternalServerError)
}
