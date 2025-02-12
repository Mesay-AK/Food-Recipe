// /main.go
package main

import (
	"food-recipe/internal/cmd/routes"
	"log"

)

func main() {
	
	router := routes.SetupRoutes();
	if err:= router.Run(":8083"); err != nil{
		log.Fatal("Server failed to start: ", err)
	}
}
