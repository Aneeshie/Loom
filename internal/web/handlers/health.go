package handlers

import (
	"net/http"

	"github.com/go-chi/render"
)

func (h *Handler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"status": "ok"}

	render.JSON(w, r, data)
}
