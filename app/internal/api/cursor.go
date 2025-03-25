package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/John-Dembaremba/pagination-technics/internal/domain/pagination"
)

type CursorBasedHttpController struct {
	Handler pagination.CursorBasedHandler
}

func NewCursorBasedHttpController(repo pagination.CursorBasedHandler) CursorBasedHttpController {
	return CursorBasedHttpController{
		Handler: repo,
	}
}

func (h CursorBasedHttpController) GetUsers(w http.ResponseWriter, r *http.Request) {
	query_params := r.URL.Query()
	cursorStr := query_params.Get("cursor")
	limitStr := query_params.Get("limit")

	var d interface{}
	cursorInt, err := strconv.Atoi(cursorStr)
	if err != nil {
		JSONResponse(w, http.StatusBadRequest, d, "invalid cursor param", "")
		return
	}

	limitInt, err := strconv.Atoi(limitStr)
	if err != nil {
		JSONResponse(w, http.StatusBadRequest, d, "invalid limit param", "")
		return
	}

	result, err := h.Handler.Retrieve(cursorInt, limitInt)
	if err != nil {
		JSONResponse(w, http.StatusInternalServerError, d, "something went wrong, please try agian", "")
		log.Printf("GetUsers Controller failed with error: %v", err)
		return
	}

	JSONResponse(w, http.StatusOK, result, "", "retrieved successfully")
	return
}
