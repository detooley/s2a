package web

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"s2a.app/db"
	"s2a.app/tsdr"
	"s2a.app/utils/structs"
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
	ids := r.URL.Query().Get("ids")
	// Query tsdr
	files := tsdr.GetFiles(ids)
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
	var idm structs.IdmPageData
	// Extract search elements
	idm.Search = r.URL.Query().Get("search")
	idm.OrderBy.Element = r.URL.Query().Get("orderby")
	idm.OrderBy.Direction = r.URL.Query().Get("direction")
	// Extract Options
	options := r.URL.Query().Get("options")
	idm.Options.Submitted = options
	for _, option := range options {
		if option == 105 {
			idm.Options.ShowId = "on"
		} // Unicode for 'i' is 105
		if option == 115 {
			idm.Options.ShowStatus = "on"
		} // Unicode for 's' is 115
		if option == 100 {
			idm.Options.ShowDate = "on"
		} // Unicode for 'd' is 100
		if option == 118 {
			idm.Options.ShowVer = "on"
		} // Unicode for 'v' is 118
		if option == 116 {
			idm.Options.ShowTm5 = "on"
		} // Unicode for 't' is 116
	}
	// Query database
	idm = db.SearchIdManual(idm)
	// Set template directory variables
	tmplDir := "web/templs/"
	tmplExt := ".html"
	// Parse template
	tmpl, err := template.New("idm" + tmplExt).ParseFiles(tmplDir + "idm" + tmplExt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Execute template
	err = tmpl.Execute(w, idm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		//	panic(err)
	}
}
