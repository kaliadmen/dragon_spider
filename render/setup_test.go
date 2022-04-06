package render

import (
	"github.com/CloudyKit/jet/v6"
	"os"
	"testing"
)

var views = jet.NewSet(jet.NewOSFileSystemLoader("./testdata/views"), jet.InDevelopmentMode())

var testRenderer = Render{
	Renderer:    "",
	RootPath:    "",
	JetTemplate: views,
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
