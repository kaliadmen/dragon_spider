package render

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type Render struct {
	Renderer    string
	JetTemplate *jet.Set
	RootPath    string
	Secure      bool
	Port        string
	ServerName  string
}

type TemplateData struct {
	IsAuthenticated bool
	IntMap          map[string]int
	StringMap       map[string]string
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
	Port            string
	ServerName      string
	Secure          bool
}

//Page renders templates using the render engine set in the Renderer
func (ren *Render) Page(w http.ResponseWriter, r *http.Request, view string, variables, data interface{}) error {
	switch strings.ToLower(ren.Renderer) {
	case "go":
		err := ren.UseGo(w, r, view, data)
		if err != nil {
			return err
		}
	case "jet":

		err := ren.UseJet(w, r, view, variables, data)
		if err != nil {
			return err
		}
	}

	return nil
}

//UseGo uses go template engine to render template pages
func (ren *Render) UseGo(w http.ResponseWriter, r *http.Request, view string, data interface{}) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/views/%s.page.tmpl", ren.RootPath, view))
	if err != nil {
		return err
	}

	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}

	err = tmpl.Execute(w, &td)
	if err != nil {
		return err
	}

	return nil
}

//UseJet uses jet engine to render template pages
func (ren *Render) UseJet(w http.ResponseWriter, r *http.Request, templateName string, variables, data interface{}) error {
	var vars jet.VarMap

	//format variables for jet
	if variables == nil {
		vars = make(jet.VarMap)
	} else {
		vars = variables.(jet.VarMap)
	}

	//format template data
	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}

	t, err := ren.JetTemplate.GetTemplate(fmt.Sprintf("%s.jet", templateName))
	if err != nil {
		log.Println(err)
		return err
	}

	if err = t.Execute(w, vars, td); err != nil {
		log.Println(err)
		return err
	}

	return nil

}
