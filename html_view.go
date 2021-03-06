package views

import tmpl "html/template"
import z "github.com/zaiuz/zaiuz"

// Represents a single html/template template. Encapsulate the template pathname from
// controller action code and allows further subviews based on this view.
type HtmlView struct {
	template  *tmpl.Template
	filenames []string
}

// Creates a new html view from the specified template name which should be a
// html/template-compatible html template file.
func NewHtmlView(filenames ...string) View {
	if len(filenames) < 1 {
		panic("needs at least 1 filename.")
	}

	view := &HtmlView{nil, filenames}
	registerAutoReparseView(view)
	return view.ReparseTemplate()
}

func (view *HtmlView) Filenames() []string {
	return view.filenames
}

func (view *HtmlView) ReparseTemplate() View {
	t, e := tmpl.ParseFiles(view.filenames...)
	noError(e) // better to failfast here since views are pre-loaded at startup.

	view.template = t
	return view
}

// Creates a subview from the receiving view. Subview templates contains all templates
// defined in the parent view.
func (view *HtmlView) Subview(filenames ...string) View {
	return NewHtmlView(append(view.filenames, filenames...)...)
}

// Renders the view to the response in the supplied Context with the given view data
// context.
func (view *HtmlView) Render(c *z.Context, data interface{}) error {
	// TODO: Configurable/overridable content type support
	w := c.ResponseWriter
	w.Header()["Content-Type"] = []string{"text/html"}
	return view.template.ExecuteTemplate(w, RootTemplateName, data)
}
