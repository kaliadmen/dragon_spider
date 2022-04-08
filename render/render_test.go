package render

import (
	"github.com/CloudyKit/jet/v6"
	"net/http"
	"net/http/httptest"
	"testing"
)

var pageData = []struct {
	name          string
	renderer      string
	template      string
	errorExpected bool
	errorMessage  string
}{
	{"using_go_engine", "go", "home", false, "error rendering go template"},
	{"using_go_engine_no_template", "go", "no-file", true, "expected error rendering non-existent template"},
	{"using_jet_engine", "jet", "home", false, "error rendering jet template"},
	{"using_jet_engine_no_template", "jet", "no-file", true, "expected error rendering non-existent template"},
	{"using_no_engine", "", "home", true, "expected error using invalid/non-existent renderer"},
}

func TestRender_Page(t *testing.T) {
	//loop over tests
	for _, test := range pageData {
		//responseWriter and request
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "/url", nil)
		if err != nil {
			t.Error(err)
		}

		testRenderer.Renderer = test.renderer
		testRenderer.RootPath = "./testdata"

		err = testRenderer.Page(w, r, test.template, nil, nil)
		if test.errorExpected {
			if err == nil {
				t.Errorf("%s: %s:", test.name, test.errorMessage)
			}

		} else {
			if err != nil {
				t.Errorf("%s: %s: %s", test.name, test.errorMessage, err.Error())
			}
		}

	}

}

func TestRender_UseGo(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/url", nil)
	if err != nil {
		t.Error(err)
	}

	testRenderer.Renderer = "go"
	testRenderer.RootPath = "./testdata"

	err = testRenderer.Page(w, r, "home", nil, nil)
	if err != nil {
		t.Error("Error rendering page", err)
	}

	td := &TemplateData{
		IsAuthenticated: false,
		IntMap:          nil,
		StringMap:       nil,
		FloatMap:        nil,
		Data:            nil,
		CSRFToken:       "",
		Port:            "",
		ServerName:      "",
		Secure:          false,
	}
	err = testRenderer.Page(w, r, "home", nil, td)
	if err != nil {
		t.Error("Error template data has nil value", err)
	}
}

func TestRender_UseJet(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/url", nil)
	if err != nil {
		t.Error(err)
	}

	testRenderer.Renderer = "jet"

	err = testRenderer.Page(w, r, "home", nil, nil)
	if err != nil {
		t.Error("Error rendering using jet", err)
	}

	vars := make(jet.VarMap)
	err = testRenderer.Page(w, r, "home", vars, nil)
	if err != nil {
		t.Error("Error Jet VarMap has nil value", err)
	}

	td := &TemplateData{
		IsAuthenticated: false,
		IntMap:          nil,
		StringMap:       nil,
		FloatMap:        nil,
		Data:            nil,
		CSRFToken:       "",
		Port:            "",
		ServerName:      "",
		Secure:          false,
	}
	err = testRenderer.Page(w, r, "home", nil, td)
	if err != nil {
		t.Error("Error template data has nil value", err)
	}
}
