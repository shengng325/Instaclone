package views

import (
	"html/template"
	"path/filepath"
)

var (
	LayoutDir   string = "views/layout/"
	TemplateExt string = ".gohtml"
)

type View struct {
	Templates *template.Template
	Layout    string
}

func NewView(layout string, files ...string) *View {
	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	newView := &View{
		Templates: t,
		Layout:    layout,
	}
	return newView
}

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}
