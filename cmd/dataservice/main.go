package main

import (
	"fmt"
	"os"

	"github.com/pulsone21/threattrack/dataservice/api"
	"github.com/pulsone21/threattrack/dataservice/storage"
)

func main() {
	fmt.Println("Setting up ThreatTrack Data-Service")
	listenAddr := os.Getenv("BACKEND_PORT")
	addr := os.Getenv("MYSQL_ADDR")
	usr := os.Getenv("MYSQL_USER")
	pw := os.Getenv("MYSQL_PW")
	port := os.Getenv("MYSQL_PORT")
	db := os.Getenv("MYSQL_DBNAME")

	dbConf := storage.DBConfig{
		Address:  addr,
		Port:     port,
		User:     usr,
		Password: pw,
		Database: db,
	}

	store := storage.NewMySqlStorage(dbConf)
	server := api.NewServer(listenAddr, store)
	fmt.Printf("Running Data-Service on http://localhost:%s\n", listenAddr)
	server.Run()

}
