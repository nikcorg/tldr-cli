package archive

import (
	"fmt"
	"io"
	"text/template"
)

// RenderToMarkdown generates Markdown output from a collection of Archive models
func RenderToMarkdown(out io.Writer, archives []Archive) {
	funcMap := template.FuncMap{
		"IsSimple":  func(x RelatedEntry) bool { return x.IsSimple() },
		"IsConcise": func(x RelatedEntry) bool { return x.IsConcise() },
	}

	if err := template.Must(
		template.New("").Funcs(funcMap).ParseFiles("./templates/layout.tmpl"),
	).ExecuteTemplate(
		out,
		"Layout",
		struct{ Archives []Archive }{archives},
	); err != nil {
		panic(fmt.Errorf("Error render: %+v", err))
	}
}
