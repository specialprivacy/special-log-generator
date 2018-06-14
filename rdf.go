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
	case "splog":
		return "http://www.specialprivacy.eu/langs/splog#" + attr
	case "dct":
		return "http://purl.org/dc/terms/" + attr
	case "prov":
		return "http://www.w3.org/ns/prov#" + attr
	case "skos":
		return "http://www.w3.org/2004/02/skos/core#" + attr
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
	tmpl := "{{$logID := randomUUID}}{{$contentId := randomUUID}}" +
		"{{if .Process}}<http://example.com/logs/{{.Process}}><http://www.w3.org/1999/02/22-rdf-syntax-ns#type><http://www.specialprivacy.eu/langs/splog#Log>;" +
		"<http://www.w3.org/ns/prov#wasAttributedTo><http://example.com/applications/{{.Process}}>;" +
		"<http://www.specialprivacy.eu/langs/splog#logEntry><http://example.com/logEntries/{{$logID}}>.{{end}}" +
		"<http://example.com/logEntries/{{$logID}}><http://www.w3.org/1999/02/22-rdf-syntax-ns#type><http://www.specialprivacy.eu/langs/splog#LogEntry>" +
		"{{if .Timestamp}};<http://www.specialprivacy.eu/langs/splog#transactionTime>\"{{toISOTime .Timestamp}}\"^^<http://www.w3.org/2001/XMLSchema#dateTime>{{end}}" +
		"{{if .UserID}};<http://www.specialprivacy.eu/langs/splog#dataSubject><http://www.example.com/users/{{.UserID}}>{{end}}" +
		";<http://www.specialprivacy.eu/langs/splog#logEntryContent><http://example.com/logEntryContents/{{$contentId}}>." +
		"<http://example.com/logEntryContents/{{$contentId}}><http://www.w3.org/1999/02/22-rdf-syntax-ns#type><http://www.specialprivacy.eu/langs/splog#LogEntryContent>" +
		"{{if .Purpose}};<http://www.specialprivacy.eu/langs/usage-policy#hasPurpose><{{.Purpose}}>{{end}}" +
		"{{if .Processing}};<http://www.specialprivacy.eu/langs/usage=policy#hasProcessing><{{.Processing}}>{{end}}" +
		"{{if .Storage}};<http://www.specialprivacy.eu/langs/usage-policy#hasStorage><{{.Storage}}>{{end}}" +
		"{{if .Recipient}};<http://www.specialprivacy.eu/langs/usage-policy#hasRecipient><{{.Recipient}}>{{end}}" +
		"{{if .Data}}{{range .Data}};<http://www.specialprivacy.eu/langs/usage-policy#hasData><{{.}}>{{end}}{{end}}."
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
			"{{if .Timestamp}};<http://purl.org/dc/terms/created>\"{{toISOTime .Timestamp}}\"^^<http://www.w3.org/2001/XMLSchema#dateTime>{{end}}" +
			"{{if .UserID}};<http://www.specialprivacy.eu/langs/usage-policy#hasDataSubject><http://www.example.com/users/{{.UserID}}>{{end}}" +
			"{{range .SimplePolicies}}" +
			";<http://www.specialprivacy.eu/vocabs/policy#simplePolicy>[" +
			"{{if .Purpose}}<http://www.specialprivacy.eu/langs/usage-policy#hasPurpose><{{.Purpose}}>{{end}}" +
			"{{if .Processing}};<http://www.specialprivacy.eu/langs/usage=policy#hasProcessing><{{.Processing}}>{{end}}" +
			"{{if .Storage}};<http://www.specialprivacy.eu/langs/usage-policy#hasStorage><{{.Storage}}>{{end}}" +
			"{{if .Recipient}};<http://www.specialprivacy.eu/langs/usage-policy#hasRecipient><{{.Recipient}}>{{end}}" +
			"{{if .Data}};<http://www.specialprivacy.eu/langs/usage-policy#hasData><{{.Data}}>{{end}}" +
			"]" +
			"{{end}}."
	output, _ := template.New("ttl-template").Funcs(funcMap).Parse(tmpl)
	return output
}
