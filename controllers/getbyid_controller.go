package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"log"
	"go-crud-generator/models"
	"go-crud-generator/validators"
)

type GetByIDController struct {
	repo   *models.DynamicRepository
	schema *models.Schema
}

func NewGetByIDController(repo *models.DynamicRepository, schema *models.Schema) *GetByIDController {
	return &GetByIDController{repo: repo, schema: schema}
}

func (gc *GetByIDController) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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

	data, err := gc.repo.FindByID(idInt)

	validators.FormatSingleDataBySchema(gc.schema, data)

	if err != nil {
		log.Printf("Erro ao buscar por ID: %v", err)
		http.Error(w, "Registro não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
