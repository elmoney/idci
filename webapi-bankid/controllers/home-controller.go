package controllers

import (
	"fmt"
	"net/http"
	"html/template"
	"github.com/idci/webapi-bankid/models"
	"log"
)

// HomeIndex главная страница
func HomeIndex(w http.ResponseWriter, r *http.Request,config *models.Config) {
	log.Println(r)
	t, error := template.ParseFiles("templates/index.html")
	if (error != nil) {
		fmt.Fprintf(w, error.Error())
	}
	t.ExecuteTemplate(w,"index",nil)
}

