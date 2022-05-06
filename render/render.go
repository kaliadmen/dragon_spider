package render

import (
	"errors"
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type Render struct {
	Renderer    string
	JetTemplate *jet.Set
	Session     *scs.SessionManager
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
		return ren.UseGo(w, r, view, data)
	case "jet":

		return ren.UseJet(w, r, view, variables, data)
	default:

	}

	return errors.New("no rendering engine found")
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

	td = ren.AddDefaultData(td, r)

	tmpl, err := ren.JetTemplate.GetTemplate(fmt.Sprintf("%s.jet", templateName))
	if err != nil {
		return err
	}

	if err = tmpl.Execute(w, vars, td); err != nil {
		log.Println(err)
		return err
	}

	return nil

}

func (ren *Render) AddDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	//check if user is authenticated
	if ren.Session.Exists(r.Context(), "userId") {
		td.IsAuthenticated = true
	}

	td.Secure = ren.Secure
	td.ServerName = ren.ServerName
	td.CSRFToken = nosurf.Token(r)
	td.Port = ren.Port

	return td
}
