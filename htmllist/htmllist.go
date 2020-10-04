package htmllist

import (
	"html/template"
	"io"
)

//Templist list all data in html template
func Templist(w io.Writer, data interface{}, tpl string, templateName string) error {
	t := template.New(templateName)
	htmllist := template.Must(t.Parse(tpl))
	err := htmllist.Execute(w, data)
	return err
}
