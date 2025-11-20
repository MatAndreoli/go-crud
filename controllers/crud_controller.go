package controllers

import (
	"html/template"
	"net/http"
	"go-crud-generator/models"
)

const defaultPageLimit = 10

type CRUDController struct {
	listController    *ListController
	createController  *CreateController
	updateController  *UpdateController
	deleteController  *DeleteController
	getByIDController *GetByIDController
}

func NewCRUDController(repo *models.DynamicRepository, schema *models.Schema, tmpl *template.Template) *CRUDController {
	return &CRUDController{
		listController:    NewListController(repo, schema, tmpl),
		createController:  NewCreateController(repo, schema, tmpl),
		updateController:  NewUpdateController(repo, schema, tmpl),
		deleteController:  NewDeleteController(repo),
		getByIDController: NewGetByIDController(repo, schema),
	}
}

func (c *CRUDController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", c.listController.HandleList)
	mux.HandleFunc("/create", c.createController.HandleCreate)
	mux.HandleFunc("/update", c.updateController.HandleUpdate)
	mux.HandleFunc("/delete", c.deleteController.HandleDelete)
	mux.HandleFunc("/get", c.getByIDController.HandleGetByID)
}
	