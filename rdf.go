package main

import (
	"fmt"
	"strings"
	"text/template"
	"time"
)

func expandPrefix(term string) string {
	splits := strings.SplitN(term, ":", 2)
	if len(splits) != 2 {
		return term
	}
	prefix := splits[0]
	attr := splits[1]
	switch prefix {
	case "spl":
		return "http://www.specialprivacy.eu/langs/usage-policy#" + attr
	case "svpu":
		return "http://www.specialprivacy.eu/vocabs/purposes#" + attr
	case "svpr":
		return "http://www.specialprivacy.eu/vocabs/processing#" + attr
	case "svr":
		return "http://www.specialprivacy.eu/vocabs/recipients#" + attr
	case "svl":
		return "http://www.specialprivacy.eu/vocabs/locations#" + attr
	case "svd":
		return "http://www.specialprivacy.eu/vocabs/data#" + attr
	default:
		return term
	}
}

func getLogTTLTemplate() *template.Template {
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
		"<http://www.specialprivacy.eu/langs/usage-policy#hasDataSubject><http://www.example.com/users/{{.UserID}}>;" +
		"<http://www.specialprivacy.eu/langs/usage-policy#hasProcessing><http://www.specialprivacy.eu/vocabs/processing#{{.Process}}>;" +
		"{{range .Attributes}}<http://www.specialprivacy.eu/langs/usage-policy#hasData><http://www.specialprivacy.eu/vocabs/data#{{.}}>;{{end}}" +
		"<http://purl.org/dc/terms/created>\"{{toISOTime .Timestamp}}\"^^<http://www.w3.org/2001/XMLSchema#dateTime>."
	output, _ := template.New("ttl-template").Funcs(funcMap).Parse(tmpl)
	return output
}

func getConsentTTLTemplate() *template.Template {
	funcMap := template.FuncMap{
		"randomUUID": randomUUID,
		"toISOTime": func(t int64) string {
			output, _ := time.Unix(0, t*int64(time.Millisecond)).MarshalText()
			return fmt.Sprintf("%s", output)
		},
	}
	// TODO: either use #hasPolicy or #hasDataSubject to link policies to a data subject (keeping both until feedback from stakeholders is received)
	tmpl :=
		"<http://www.example.com/users/{{.UserID}}><http://www.specialprivacy.eu/langs/usage-policy#hasPolicy><http://www.example.com/policy/{{.ConsentID}}>." +
			"<http://www.example.com/policy/{{.ConsentID}}><http://www.w3.org/1999/02/22-rdf-syntax-ns#type><http://www.specialprivacy.eu/vocabs/policy#Consent>" +
			"{{if .Purpose}};<http://www.specialprivacy.eu/langs/usage-policy#hasPurpose><{{.Purpose}}>{{end}}" +
			"{{if .Processing}};<http://www.specialprivacy.eu/langs/usage=policy#hasProcessing><{{.Processing}}>{{end}}" +
			"{{if .Storage}};<http://www.specialprivacy.eu/langs/usage-policy#hasStorage><{{.Storage}}>{{end}}" +
			"{{if .Recipient}};<http://www.specialprivacy.eu/langs/usage-policy#hasRecipient><{{.Recipient}}>{{end}}" +
			"{{if .UserID}};<http://www.specialprivacy.eu/langs/usage-policy#hasDataSubject><http://www.example.com/users/{{.UserID}}>{{end}}" +
			"{{if .Data}};<http://www.specialprivacy.eu/langs/usage-policy#hasData><{{.Data}}>{{end}}" +
			"{{if .Timestamp}};<http://purl.org/dc/terms/created>\"{{toISOTime .Timestamp}}\"^^<http://www.w3.org/2001/XMLSchema#dateTime>{{end}}."
	output, _ := template.New("ttl-template").Funcs(funcMap).Parse(tmpl)
	return output
}
