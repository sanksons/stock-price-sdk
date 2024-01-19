package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sanksons/stock-price-sdk/pkg"
)

const portNum string = ":8080"

// Handler functions.
func Home(w http.ResponseWriter, r *http.Request) {
	stock := r.URL.Query().Get("stock")
	stockNav, err := pkg.GetStockNav(stock)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))

		return
	}
	data, err := json.Marshal(stockNav)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(data)
}

func main() {
	log.Println("Starting our simple http server.")

	// Registering our handler functions, and creating paths.
	http.HandleFunc("/home", Home)
	http.HandleFunc("/", Home)

	log.Println("Started on port", portNum)
	fmt.Println("To close connection CTRL+C :-)")

	// Spinning up the server.
	err := http.ListenAndServe(portNum, nil)
	if err != nil {
		log.Fatal(err)
	}
}
