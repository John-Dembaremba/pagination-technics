package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/John-Dembaremba/pagination-technics/internal/domain/pagination"
	"github.com/John-Dembaremba/pagination-technics/pkg"
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
	// tracer span instance
	tracerHander := pkg.TracerConfigHandler{}
	ctx, span := tracerHander.TracerSpan(r.Context(), "cursor-httpController", "controller: get-users")
	defer span.End()

	// query params handling
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

	// domain layer
	result, err := h.Handler.Retrieve(ctx, cursorInt, limitInt)
	if err != nil {
		JSONResponse(w, http.StatusInternalServerError, d, "something went wrong, please try agian", "")
		span.RecordError(err) // Record error in span
		log.Printf("GetUsers Controller failed with error: %v", err)
		return
	}

	JSONResponse(w, http.StatusOK, result, "", "retrieved successfully")
	return
}
