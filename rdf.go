package main

import (
	"fmt"
	"text/template"
	"time"
)

func getTtlTemplate() *template.Template {
	funcMap := template.FuncMap{
		"randomUUID": randomUUID,
		"toISOTime": func(t int64) string {
			output, _ := time.Unix(0, t*int64(time.Millisecond)).MarshalText()
			return fmt.Sprintf("%s", output)
		},
	}
	tmpl := "<http://www.specialprivacy.eu/log/{{randomUUID}}><http://www.w3.org/1999/02/22-rdf-syntax-ns#type><http://www.specialprivacy.eu/vocabs/logs#log>;" +
		"<http://www.specialprivacy.eu/langs/usage-policy#hasPurpose><http://www.specialprivacy.eu/vocabs/purposes#{{.Purpose}}>;" +
		"<http://www.specialprivacy.eu/langs/usage-policy#hasStorage><http://www.specialprivacy.eu/vocabs/locations#{{.Location}}>;" +
		"<http://www.specialprivacy.eu/langs/usage-policy#hasDataSubject><http://www.example.com/users/{{.UserId}}>;" +
		"<http://www.specialprivacy.eu/langs/usage-policy#hasProcessing><http://www.specialprivacy.eu/vocabs/processing#{{.Process}}>;" +
		"{{range .Attributes}}<http://www.specialprivacy.eu/langs/usage-policy#hasData><http://www.specialprivacy.eu/vocabs/data#{{.}}>;{{end}}" +
		"<http://purl.org/dc/terms/created>\"{{toISOTime .Timestamp}}\"^^<http://www.w3.org/2001/XMLSchema#dateTime>."
	output, _ := template.New("ttl-template").Funcs(funcMap).Parse(tmpl)
	return output
}
