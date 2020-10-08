package views

import (
	"bytes"
	// "fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"lenslocked.com/context"
)

var (
	TemplateDir string = "views/"
	LayoutDir   string = "views/layout/"
	TemplateExt string = ".gohtml"
)

type View struct {
	Template *template.Template
	Layout   string
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}

func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	var vd Data
	switch d := data.(type) {
	case Data:
		vd = d
	default:
		vd = Data{
			Yield: data,
		}
	}
	vd.User = context.User(r.Context())
	var buf bytes.Buffer
	err := v.Template.ExecuteTemplate(&buf, v.Layout, vd)
	if err != nil {
		http.Error(w, "Something went wrong. If the problem epersists, please email support@lenslocked.com", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
	// return v.Template.ExecuteTemplate(w, v.Layout, data)
}

func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	newView := &View{
		Template: t,
		Layout:   layout,
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

func addTemplatePath(files []string) {
	for i, file := range files {
		files[i] = TemplateDir + file
	}
}

func addTemplateExt(files []string) {
	for i, file := range files {
		files[i] = file + TemplateExt
	}
}
