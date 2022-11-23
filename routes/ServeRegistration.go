package routes

import (
	"html/template"
	"log"
	"net/http"
)

func ServeRegistration(origin string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Origin string
		}{origin}

		temp, err := template.ParseFiles("../../html/registration.gohtml")
		if err != nil {
			log.Fatal("Unable to parse registration page template:", err)
		}

		err = temp.Execute(w, data)
		if err != nil {
			log.Println("Error while parsing registration template:", err)
		}
	}
}
