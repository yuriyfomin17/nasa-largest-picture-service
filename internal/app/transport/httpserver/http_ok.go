package httpserver

import (
	"encoding/json"
	"net/http"
)

func RespondOk(data any, w *http.ResponseWriter) {
	(*w).Header().Add("Content-Type", "application/json; charset=utf-8")
	(*w).WriteHeader(http.StatusOK)
	_ = json.NewEncoder(*w).Encode(data)
}
