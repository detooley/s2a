package web

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"s2a.app/db"
	"s2a.app/tsdr"
	"s2a.app/utils/text"
)

// Root
func GetRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /request\n")
	var page string
	if r.URL.Path != "/" {
		page = "404"
	} else {
		page = "home"
	}
	// Set template directory variables
	tmplDir := "web/templs/"
	tmplExt := ".html"
	// Parse template
	tmpl, err := template.New(page + tmplExt).ParseFiles(tmplDir + page + tmplExt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Execute template
	err = tmpl.Execute(w, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		//	panic(err)
	}
}

// Generic
func GetGen(w http.ResponseWriter, r *http.Request) {
	page := text.ExtractBetween(r.URL.String(), "/", "?")
	fmt.Print("got /" + page + "\n")
	// Extract serial numbers
	sns := r.URL.Query().Get("sns")
	// Query tsdr
	files := tsdr.GetFiles(sns)
	//fmt.Print(files)
	// Set template functions
	funcMap := template.FuncMap{
		"dec": func(i int) int { return i - 1 },
		"slz": func(i string) string { return strings.TrimLeft(i, "0") },
		"brc": text.RemoveBrackets,
	}
	// Set template directory variables
	tmplDir := "web/templs/"
	tmplExt := ".html"
	// Parse template
	tmpl, err := template.New(page + tmplExt).Funcs(funcMap).ParseFiles(tmplDir + page + tmplExt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Execute template
	err = tmpl.Execute(w, files)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		//	panic(err)
	}
}

// Db Creation
func GetDb(w http.ResponseWriter, r *http.Request) {
	fmt.Print("got /db\n")
	db.CreateIdManual()
	// Set template directory variables
	tmplDir := "web/templs/"
	tmplExt := ".html"
	// Parse template
	tmpl, err := template.New("db" + tmplExt).ParseFiles(tmplDir + "db" + tmplExt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Execute template
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		//	panic(err)
	}
}

// Search Id Manual
func GetIdm(w http.ResponseWriter, r *http.Request) {
	fmt.Print("got /idm\n")
	// Extract search term
	search := r.URL.Query().Get("search")
	entries := db.SearchIdManual(search)
	// Set template directory variables
	tmplDir := "web/templs/"
	tmplExt := ".html"
	// Parse template
	tmpl, err := template.New("idm" + tmplExt).ParseFiles(tmplDir + "idm" + tmplExt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Execute template
	err = tmpl.Execute(w, entries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		//	panic(err)
	}
}
