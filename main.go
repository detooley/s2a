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
	// Serve webpages
	http.HandleFunc("/", web.GetRoot)
	http.HandleFunc("/cites", web.GetGen)
	http.HandleFunc("/marks", web.GetGen)
	http.HandleFunc("/reds", web.GetGen)
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
