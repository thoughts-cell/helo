package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]
		data, err := query(city)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	})

	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

//define a way to populate the structure
//takes a string representing the city, and returns a weatherData struct and an error
func query(city string) (weatherData, error) {
	//return weatherData type and error type
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=YOUR_API_KEY&q=" + city)

	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()
	var d weatherData
	//the json.NewDecoder leverages an elegant feature of Go, which are interfaces.
	//The Decoder doesn’t take a concrete HTTP response body; rather, it takes an io.Reader interface,
	// which the http.Response.Body happens to satisfy.
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}
	return d, nil //Finally, if the decode succeeds, we return the weatherData to the caller,
	//with a nil error to indicate success.
}

type weatherData struct {
	//only capture the information we care about
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}
