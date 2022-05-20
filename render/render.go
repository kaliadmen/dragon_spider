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
	Error           string
	Flash           string
}

//Page renders templates using the render engine set in the Renderer
func (r *Render) Page(w http.ResponseWriter, req *http.Request, view string, variables, data interface{}) error {
	switch strings.ToLower(r.Renderer) {
	case "go":
		return r.UseGo(w, req, view, data)
	case "jet":

		return r.UseJet(w, req, view, variables, data)
	default:

	}

	return errors.New("no rendering engine found")
}

//UseGo uses go template engine to render template pages
func (r *Render) UseGo(w http.ResponseWriter, req *http.Request, view string, data interface{}) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/views/%s.page.tmpl", r.RootPath, view))
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
func (r *Render) UseJet(w http.ResponseWriter, req *http.Request, templateName string, variables, data interface{}) error {
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

	td = r.AddDefaultData(td, req)

	tmpl, err := r.JetTemplate.GetTemplate(fmt.Sprintf("%s.jet", templateName))
	if err != nil {
		return err
	}

	if err = tmpl.Execute(w, vars, td); err != nil {
		log.Println(err)
		return err
	}

	return nil

}

func (r *Render) AddDefaultData(td *TemplateData, req *http.Request) *TemplateData {
	//check if user is authenticated
	if r.Session.Exists(req.Context(), "userId") {
		td.IsAuthenticated = true
	}

	td.Secure = r.Secure
	td.ServerName = r.ServerName
	td.CSRFToken = nosurf.Token(req)
	td.Port = r.Port
	td.Error = r.Session.PopString(req.Context(), "error")
	td.Flash = r.Session.PopString(req.Context(), "flash")

	return td
}
