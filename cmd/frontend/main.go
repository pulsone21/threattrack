package main

import (
	"fmt"
	"os"

	"github.com/pulsone21/threattrack/frontend/gui"
)

func main() {
	fmt.Println("Setting up FrontEnd")
	port := os.Getenv("FRONTEND_PORT")
	backendAdress := fmt.Sprintf("%s:%s", os.Getenv("BACKEND_ADRESS"), os.Getenv("BACKEND_PORT"))
	fmt.Println("BackendAdress: ", backendAdress)
	server := gui.CreateServer(port, backendAdress)
	fmt.Println("trying to start frontEnd Service")
	fmt.Println(server.Port)
	server.Run()

}
