package util

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mxrcury/rootgo/types"
)

func DecodeBody(ctx context.Context) *json.Decoder {
	bodyJSON := string(ctx.Value("body").([]byte))
	bodyReader := strings.NewReader(bodyJSON)
	return json.NewDecoder(bodyReader)
}

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
