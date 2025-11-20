package controllers

import (
	"net/http"
	"strconv"
	"log"
	"go-crud-generator/models"
)

type DeleteController struct {
	repo *models.DynamicRepository
}

func NewDeleteController(repo *models.DynamicRepository) *DeleteController {
	return &DeleteController{repo: repo}
}

func (dc *DeleteController) HandleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID ausente", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := dc.repo.Delete(idInt); err != nil {
		log.Printf("Erro ao deletar registro: %v", err)
		http.Error(w, "Erro ao deletar registro", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
