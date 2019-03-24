package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"

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

const apiKey = "7d2e58d76328a93bb73cb74a167a58d7"

func main() {

	templates := template.Must(template.ParseFiles("templates/index.html"))
	// db, _ := sql.Open("sqlite3", "dev.db")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// p := person{Name: "Aditya"}
		// name := r.FormValue("name")
		// if name != "" {
		// 	p.Name = name
		// }
		// fmt.Print(name)
		// p.DBStatus = db.Ping() == nil
		err := templates.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		var results []Results
		var err error

		if results, err = search(r.FormValue("search")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		encoder := json.NewEncoder(w)
		err = encoder.Encode(results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Println(http.ListenAndServe(":8080", nil))
}

func search(query string) ([]Results, error) {
	var resp *http.Response
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://developers.zomato.com/api/v2.1/locations?query="+url.QueryEscape(query), nil)
	if err != nil {
		return []Results{}, err
	}
	req.Header.Add("user-key", apiKey)
	resp, err = client.Do(req)
	if err != nil {
		return []Results{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Results{}, err
	}
	fmt.Print(string(body))
	return []Results{}, err
}
