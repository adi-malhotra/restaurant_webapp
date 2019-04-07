package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/codegangsta/negroni"
	_ "github.com/mattn/go-sqlite3"
	"github.com/yosssi/ace"
)

// SearchResults : struct to hold results
type SearchResults struct {
	Results []Restaurants `json:"restaurants,attr"`
}

// Restaurants : struct to hold a restaurant
type Restaurants struct {
	Restaurant Restaurant `json:"restaurant"`
}

// Restaurant : Attributes of a restaurant
type Restaurant struct {
	ID       string `json:"id"`
	Name     string `json:"name,attr"`
	Location struct {
		Address string `json:"address"`
		City    string `json:"city"`
	} `json:"location"`
	Cuisines   string `json:"cuisines,attr"`
	UserRating struct {
		AggregateRating interface{} `json:"aggregate_rating"`
		RatingText      string      `json:"rating_text"`
		RatingColor     string      `json:"rating_color"`
	} `json:"user_rating"`
	AverageCostForTwo int `json:"average_cost_for_two"`
}

const apiKey = "7d2e58d76328a93bb73cb74a167a58d7"

var db *sql.DB

func verifyDatabase(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if err := db.Ping(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	next(w, r)
}

func main() {
	db, _ = sql.Open("sqlite3", "dev.db")

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template, err := ace.Load("templates/index", "", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = template.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
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

	mux.HandleFunc("/restaurant/add", func(w http.ResponseWriter, r *http.Request) {
		var restaurant Restaurant
		var err error
		if restaurant, err = findRestaurant(r.FormValue("id")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		_, err = db.Exec("insert into restaurant(sno, name, location, cuisines, user_rating, cost_for_two, rest_id) values (?,?,?,?,?,?,?)",
			nil, restaurant.Name, restaurant.Location.Address, restaurant.Cuisines, restaurant.UserRating.AggregateRating, restaurant.AverageCostForTwo, restaurant.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	})

	n := negroni.Classic()
	n.Use(negroni.HandlerFunc(verifyDatabase))
	n.UseHandler(mux)
	n.Run(":8080")
}

func search(query string) ([]Restaurants, error) {
	var resp *http.Response
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://developers.zomato.com/api/v2.1/search?q="+url.QueryEscape(query)+"&count=30", nil)
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

func findRestaurant(query string) (Restaurant, error) {
	var resp *http.Response
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://developers.zomato.com/api/v2.1/restaurant?res_id="+url.QueryEscape(query), nil)
	if err != nil {
		return Restaurant{}, err
	}
	req.Header.Add("user-key", apiKey)
	resp, err = client.Do(req)
	if err != nil {
		return Restaurant{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Restaurant{}, err
	}
	// fmt.Println(string(body))
	var rest Restaurant
	if err = json.Unmarshal(body, &rest); err != nil {
		panic(err)
	}
	return rest, err
}
