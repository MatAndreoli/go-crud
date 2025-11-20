package controllers

import (
	"html/template"
	"net/http"
	"log"
	"go-crud-generator/models"
	"go-crud-generator/validators"
)

type CreateController struct {
	repo   *models.DynamicRepository
	schema *models.Schema
	tmpl   *template.Template
}

func NewCreateController(repo *models.DynamicRepository, schema *models.Schema, tmpl *template.Template) *CreateController {
	return &CreateController{repo: repo, schema: schema, tmpl: tmpl}
}

func (cc *CreateController) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erro ao parsear formulário", http.StatusBadRequest)
		return
	}

	data, validationErrors := validators.ValidateData(r.PostForm, cc.schema)

	if len(validationErrors) > 0 {
		ReloadPageWithErrors(cc.tmpl, cc.repo, cc.schema, w, r, validationErrors, r.PostForm)
		return
	}

	_, err := cc.repo.Create(data)
	if err != nil {
		log.Printf("Erro ao criar registro: %v", err)
		validationErrors["_form"] = "Erro interno ao salvar. Verifique se os dados estão corretos."
		ReloadPageWithErrors(cc.tmpl, cc.repo, cc.schema, w, r, validationErrors, r.PostForm)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}