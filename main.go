package main

import (
	"fmt"
	"image_server/config"
	"image_server/middleware"
	"net/http"
)

func main() {
	fmt.Println("--------------------Image Store Service Started--------------------")

	router := middleware.NewRouter()
	port := config.GetConfigParamString("bind_port")

	http.ListenAndServe(":"+port, router)
}
