package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.gohtml")
	})

	fmt.Println(os.Getwd())

	fmt.Println("Starting front end service in port 8081")

	log.Fatal(http.ListenAndServe(":8081", nil))
}

//go:embed templates
var templateFS embed.FS

func render(w http.ResponseWriter, t string) {
	partials := []string{
		"templates/base.layout.gohtml",
		"templates/header.partial.gohtml",
		"templates/footer.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("templates/%s", t))

	// for _, x := range partials {
	// 	templateSlice = append(templateSlice, x)
	// }
	templateSlice = append(templateSlice, partials...)

	tmpl, err := template.ParseFS(templateFS, templateSlice...)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data struct {
		BrokerURL string
	}
	data.BrokerURL = os.Getenv("BROKER_URL")
	// data.BrokerURL = "http://localhost:8080"

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
