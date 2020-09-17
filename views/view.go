package views

import "html/template"

type View struct {
	Templates *template.Template
	Layout    string
}

func NewView(layout string, files ...string) *View {
	files = append(files, "views/layout/bootstrap.gohtml", "views/layout/footer.gohtml")
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
