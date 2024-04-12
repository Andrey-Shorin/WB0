package main

import (
	"fmt"
	"main/internal/app"
	"main/internal/config"
	"net/http"
)

func main() {

	var myApp app.App
	myApp.Init(config.ReadConfig())
	defer myApp.Close()
	http.HandleFunc("/", myApp.Handler)
	fmt.Printf("App started\n")
	http.ListenAndServe(":8080", nil)

}
