package api

import (
	"net/http"
	"strconv"

	"github.com/John-Dembaremba/pagination-technics/internal/domain/pagination"
)

type LimitOffsetHttpControler struct {
	Handler pagination.LimitOffSetHandler
}

func (h LimitOffsetHttpControler) GetUsers(w http.ResponseWriter, r *http.Request) {
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

	userData, err := h.Handler.RetrieveUsers(pageInt, limitInt)
	if err != nil {
		JSONResponse(w, http.StatusInternalServerError, d, "something went wrong", "")
		return
	}

	JSONResponse(w, http.StatusOK, userData, "", "retrieved successfully")

}
