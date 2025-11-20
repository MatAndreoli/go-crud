package controllers

import (
	"html/template"
	"math"
	"net/http"
	"strconv"
	"time"
	"log"

	"go-crud-generator/models"
	"go-crud-generator/validators"
)

type ListController struct {
	repo   *models.DynamicRepository
	schema *models.Schema
	tmpl   *template.Template
}

func NewListController(repo *models.DynamicRepository, schema *models.Schema, tmpl *template.Template) *ListController {
	return &ListController{repo: repo, schema: schema, tmpl: tmpl}
}

func (lc *ListController) HandleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	pageStr := r.URL.Query().Get("page")
	dbSearch := r.URL.Query().Get("search")
	page, _ := strconv.Atoi(pageStr)
	if page <= 0 {
		page = 1
	}

	data, totalRecords, err := lc.repo.FindAll(page, 10, dbSearch)
	if err != nil {
		log.Printf("Erro ao buscar dados: %v", err)
		http.Error(w, "Erro ao buscar dados", http.StatusInternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalRecords) / 10.0))
	pagination := Pagination{
		CurrentPage:  page,
		TotalPages:   totalPages,
		TotalRecords: totalRecords,
		HasPrev:      page > 1,
		PrevPage:     page - 1,
		HasNext:      page < totalPages,
		NextPage:     page + 1,
	}

	validators.FormatDataBySchema(lc.schema, data)

	templateData := TemplateData{
		Schema:        lc.schema,
		Data:          data,
		SearchTerm:    dbSearch,
		Pagination:    pagination,
		CurrentTime:   time.Now().Unix(),
		SchemaColspan: len(lc.schema.Fields) + 1,
	}

	if err := lc.tmpl.ExecuteTemplate(w, "crud.html", templateData); err != nil {
		log.Printf("Erro ao renderizar template: %v", err)
		http.Error(w, "Erro ao renderizar página", http.StatusInternalServerError)
	}
}
