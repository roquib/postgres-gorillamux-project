package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/roquib/go-postgres/router"
)

func main() {
	r := router.Router()
	fmt.Println("Starting server on the port 9000")
	log.Fatal(http.ListenAndServe(":9000", r))
}
