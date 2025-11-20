package controllers

import (
	"html/template"
	"net/http"
	"log"
	"go-crud-generator/models"
	"go-crud-generator/validators"
)

type UpdateController struct {
	repo   *models.DynamicRepository
	schema *models.Schema
	tmpl   *template.Template
}

func NewUpdateController(repo *models.DynamicRepository, schema *models.Schema, tmpl *template.Template) *UpdateController {
	return &UpdateController{repo: repo, schema: schema, tmpl: tmpl}
}

func (uc *UpdateController) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erro ao parsear formulário", http.StatusBadRequest)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "ID ausente", http.StatusBadRequest)
		return
	}

	data, validationErrors := validators.ValidateData(r.PostForm, uc.schema)
	if len(validationErrors) > 0 {
		ReloadPageWithErrors(uc.tmpl, uc.repo, uc.schema, w, r, validationErrors, r.PostForm)
		return
	}

	var pkValue interface{}
	for _, field := range uc.schema.Fields {
		if field.PrimaryKey {
			pkValue = data[field.Name]
			break
		}
	}

	if err := uc.repo.Update(pkValue, data); err != nil {
		log.Printf("Erro ao atualizar registro: %v", err)
		validationErrors["_form"] = "Erro interno ao atualizar."
		ReloadPageWithErrors(uc.tmpl, uc.repo, uc.schema, w, r, validationErrors, r.PostForm)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}