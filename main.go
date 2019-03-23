package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type person struct {
	Name     string
	DBStatus bool
}

type Results struct {
	RestaurantName string
	Location       string
	Cuisines       string
	UserRating     float32
	AvgCost        int
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
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		results := []Results{
			Results{"Levels", "Hauz Khas Village", "Italian, North Indian, Mexican", 4.2, 1200},
			Results{"Local", "Connaught Place", "Buffet", 3.9, 1500},
			Results{"Diggin", "Chankyapuri", "Mexican, Greek, Italian, Indian", 4.4, 2000},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		encoder := json.NewEncoder(w)
		err := encoder.Encode(results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	fmt.Println(http.ListenAndServe(":8080", nil))
}
