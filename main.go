package main

import (
	"flag"
	"log"
	"os"

	"git.codesubmit.io/sfox/party-invite-ruiegv/api"
)

func main() {
	logPath := flag.String("logPath", "log.txt", "Path of log file")
	port := flag.String("port", "8081", "Listening port")
	officeLatitude := flag.Float64("latitude", 53.339428, "Latitude of office")
	officeLongitude := flag.Float64("longitude", -6.257664, "Longitude of office")

	flag.Parse()
	//init the logger with the specified path
	logWriter, err := os.OpenFile(*logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Fail to open log file: ", err.Error())
		return
	}
	log.SetOutput(logWriter)

	//Get the api instance
	api, err := api.GetApiV1(*officeLongitude, *officeLatitude)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	//Start the api
	if err := api.StartServer(":" + *port); err != nil {
		log.Fatal("Fail to start server: ", err.Error())
		return
	}
}
