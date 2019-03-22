package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type person struct {
	Name     string
	DBStatus bool
}

func main() {

	templates := template.Must(template.ParseFiles("templates/index.html"))
	db, _ := sql.Open("sqlite3", "dev.db")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := person{Name: "Aditya"}
		name := r.FormValue("name")
		if name != "" {
			p.Name = name
		}
		p.DBStatus = db.Ping() == nil
		err := templates.ExecuteTemplate(w, "index.html", p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	fmt.Println(http.ListenAndServe(":8080", nil))
}
