package main

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"s2a.app/db"
)

func main() {
	c := cron.New()

	// Add a job that runs every Thursday at 3:00 AM to update Id Manual
	_, err := c.AddFunc("0 3 * * 4", func() {
		db.CreateIdManual()
		log.Println("Id Manual Created at:", time.Now().Format(time.RFC3339))
	})
	if err != nil {
		log.Fatal(err)
	}

	// Add a job that runs every 30 secs to update Pulls
	//_, err = c.AddFunc("@every 30s", func() {
	//	results := db.SearchDb("SELECT * FROM pulls ORDER BY pull_date DESC LIMIT 1;")
	//	serNo := int(results[0]["ser_no"].(int64))
	//	for i := range 101 {
	//		serNo := serNo + i
	//		fmt.Println(serNo)
	//	}
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}

	// Start the scheduler
	c.Start()

	// Keep the main function running to allow cron jobs to execute
	select {}
}
