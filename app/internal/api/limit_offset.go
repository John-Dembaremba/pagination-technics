package api

import (
	"net/http"
	"strconv"

	"github.com/John-Dembaremba/pagination-technics/internal/domain/pagination"
	"github.com/John-Dembaremba/pagination-technics/pkg"
)

type LimitOffsetHttpControler struct {
	Handler pagination.LimitOffSetHandler
}

func NewLimitOffsetHttpControler(repo pagination.LimitOffSetHandler) LimitOffsetHttpControler {
	return LimitOffsetHttpControler{
		Handler: repo,
	}
}

func (h LimitOffsetHttpControler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// tracer span instance
	tracerHander := pkg.TracerConfigHandler{}
	ctx, span := tracerHander.TracerSpan(r.Context(), "limit-offset-httpController", "controller: get-users")
	defer span.End()

	url := r.URL.Query()

	pageStr := url.Get("page")
	limitStr := url.Get("limit")

	var d interface{}
	pageInt, err := strconv.Atoi(pageStr)
	if err != nil {
		JSONResponse(w, http.StatusBadRequest, d, "invalid page", "")
		return
	}

	limitInt, err := strconv.Atoi(limitStr)
	if err != nil {
		JSONResponse(w, http.StatusBadRequest, d, "invalid limit", "")
		return
	}

	userData, err := h.Handler.RetrieveUsers(ctx, pageInt, limitInt)
	if err != nil {
		span.RecordError(err)
		JSONResponse(w, http.StatusInternalServerError, d, "something went wrong", "")
		return
	}

	JSONResponse(w, http.StatusOK, userData, "", "retrieved successfully")

}
