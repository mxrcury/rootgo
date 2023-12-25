package util

import (
	"encoding/json"
	"net/http"

	"github.com/mxrcury/rootgo/types"
)

func WriteJSON(w http.ResponseWriter, data interface{}, status int) {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, err types.Error) {
	w.WriteHeader(err.Status)
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Connection", "close")
	json.NewEncoder(w).Encode(err)
}
