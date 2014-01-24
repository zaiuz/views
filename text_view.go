package views

import tmpl "text/template"
import z "github.com/zaiuz/zaiuz"

type TextView struct {
	template  *tmpl.Template
	filenames []string
}

func NewTextView(filenames ...string) View {
	if len(filenames) < 1 {
		panic("needs at least 1 filename.")
	}

	t, e := tmpl.ParseFiles(filenames...)
	if e != nil {
		panic(e) // better to failfast here
	}

	return &TextView{t, filenames}
}

func (view *TextView) Subview(filenames ...string) View {
	return NewTextView(append(view.filenames, filenames...)...)
}

func (view *TextView) Render(c *z.Context, data interface{}) error {
	w := c.ResponseWriter
	w.Header()["Content-Type"] = []string{"text/plain"}
	return view.template.ExecuteTemplate(w, RootTemplateName, data)
}
