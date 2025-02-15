package api

import (
	"encoding/json"
	"net/http"

	"github.com/John-Dembaremba/pagination-technics/internal/model"
)

func JSONResponse(w http.ResponseWriter, status int, result interface{}, errMsg, successMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	var response model.ResponseMeta
	if errMsg != "" {
		response.Error = errMsg
	}
	if successMsg != "" {
		response.Success = successMsg
	}

	response.Data = result

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
