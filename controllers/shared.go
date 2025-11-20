package controllers

import (
	"math"
	"net/http"
	"strconv"
	"time"
	"html/template"
	"log"
	"go-crud-generator/models"
)

type Pagination struct {
	CurrentPage  int
	TotalPages   int
	TotalRecords int
	HasPrev      bool
	HasNext      bool
	PrevPage     int
	NextPage     int
}

type TemplateData struct {
	Schema         *models.Schema
	Data           []map[string]interface{}
	Errors         map[string]string
	FormData       map[string]string
	SearchTerm     string
	Pagination     Pagination
	CurrentTime    int64
	SuccessMessage string
	SchemaColspan  int
}

func ReloadPageWithErrors(tmpl *template.Template, repo *models.DynamicRepository, schema *models.Schema, w http.ResponseWriter, r *http.Request, errors map[string]string, formData map[string][]string) {
	pageStr := r.URL.Query().Get("page")
	search := r.URL.Query().Get("search")
	page, _ := strconv.Atoi(pageStr)
	if page <= 0 {
		page = 1
	}

	data, totalRecords, err := repo.FindAll(page, 10, search)
	if err != nil {
		log.Printf("Erro ao buscar dados: %v", err)
		http.Error(w, "Erro ao buscar dados", http.StatusInternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalRecords) / 10.0))
	pagination := Pagination{
		CurrentPage: page,
		TotalPages:  totalPages,
		HasPrev:     page > 1,
		PrevPage:    page - 1,
		HasNext:     page < totalPages,
		NextPage:    page + 1,
	}

	simpleFormData := make(map[string]string)
	for k, v := range formData {
		if len(v) > 0 {
			simpleFormData[k] = v[0]
		}
	}

	templateData := TemplateData{
		Schema:        schema,
		Data:          data,
		SearchTerm:    search,
		Pagination:    pagination,
		Errors:        errors,
		FormData:      simpleFormData,
		CurrentTime:   time.Now().Unix(),
		SchemaColspan: len(schema.Fields) + 1,
	}

	w.WriteHeader(http.StatusBadRequest)
	err2 := tmpl.ExecuteTemplate(w, "crud.html", templateData)
	if err2 != nil {
		log.Printf("Erro ao renderizar template: %v", err2)
		http.Error(w, "Erro ao renderizar página", http.StatusInternalServerError)
	}
}
