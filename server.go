package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Page struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		rawurl := r.FormValue("url")
		if rawurl == "" {
			http.Error(w, "url is not specified", http.StatusBadRequest)
			return
		}

		enc := json.NewEncoder(w)
		var p Page
		//	var list []

		p.Title = "sample"
		p.Description = "description"

		if err := enc.Encode(&p); err != nil {
			fmt.Println("failed")
		}

	})
	http.ListenAndServe(":8080", nil)
}
