package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Aneeshie/loom/internal/web/service"
	"github.com/Aneeshie/loom/internal/web/types"
)

type Handler struct {
	queryService *service.QueryService
}

func NewHandler(queryService *service.QueryService) *Handler {
	return &Handler{
		queryService: queryService,
	}
}

func (h *Handler) HandleQuery(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	//get the query from the request
	var req types.QueryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.queryService.Query(r.Context(), req.Query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// using service get the response

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode resposne", http.StatusInternalServerError)
		return
	}

	// return it thats it
}
