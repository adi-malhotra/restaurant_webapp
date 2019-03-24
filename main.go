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

// Restaurant : Attributes of a restaurant
type Restaurant struct {
	Name     string `json:"name,attr"`
	Location struct {
		Address string `json:"address"`
		City    string `json:"city"`
	} `json:"location"`
	Cuisines   string `json:"cuisines,attr"`
	UserRating struct {
		AggregateRating string `json:"aggregate_rating"`
		RatingText      string `json:"rating_text"`
		RatingColor     string `json:"rating_color"`
	} `json:"user_rating"`
	AverageCostForTwo int `json:"average_cost_for_two"`
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
		var results []Restaurants
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

// SearchResults : struct to hold results
type SearchResults struct {
	Results []Restaurants `json:"restaurants,attr"`
}

// Restaurants : struct to hold a restaurant
type Restaurants struct {
	Restaurant Restaurant `json:"restaurant"`
}

func search(query string) ([]Restaurants, error) {
	var resp *http.Response
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://developers.zomato.com/api/v2.1/search?q="+url.QueryEscape(query)+"&count=20", nil)
	if err != nil {
		return []Restaurants{}, err
	}
	req.Header.Add("user-key", apiKey)
	resp, err = client.Do(req)
	if err != nil {
		return []Restaurants{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Restaurants{}, err
	}
	var s SearchResults
	if err = json.Unmarshal(body, &s); err != nil {
		panic(err)
	}
	// fmt.Println(s.Results)
	return s.Results, err
}
