package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"s2a.app/web"
)

// MAIN

func main() {
	// Start cron jobs
	//cidm.StartIdmCron()
	// Serve static files
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("web/css"))))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/favicon/favicon.ico")
	})
	// Serve webpages
	http.HandleFunc("/", web.GetRoot)
	http.HandleFunc("/cites", web.GetCites)
	http.HandleFunc("/marks", web.GetTsdr)
	http.HandleFunc("/reds", web.GetTsdr)
	http.HandleFunc("/idm", web.GetIdm)
	http.HandleFunc("/db", web.GetDb)

	// Server
	err := http.ListenAndServe(":3000", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
