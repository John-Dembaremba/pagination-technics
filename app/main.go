package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/John-Dembaremba/pagination-technics/pkg"
)

func main() {
	log.Println("Setting Env Variables ...")
	env := pkg.NewEnv()

	log.Printf("Starting Server on port: %v\n", env.ServerPort)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello Paginators are ready")
	})

	http.ListenAndServe(":3025", mux)
}
