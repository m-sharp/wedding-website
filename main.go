package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type RenderContext struct {
	TargetYear int
}

func main() {
	pageContext := &RenderContext{
		TargetYear: 2023,
	}

	templates := template.Must(template.ParseFiles("templates/home.html"))

	// Handle static asset requests
	http.Handle(
		"/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	// Handle other requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := templates.ExecuteTemplate(w, "home.html", pageContext); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	//Start the web server, set the port to listen to 8080. Without a path it assumes localhost
	//Print any errors from starting the webserver using fmt
	println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		println(fmt.Errorf("stopped listening with error: %w", err).Error())
	}
}
