package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/laurentlouk/weather-color/utils"
)

const (
	// Api Key
	key     = "32bac24471194d60dfed6d0b3370f4fa"
	version = "2.5"
	// TempRedRange : From -20°C to 50°C, 50 + 20 + 1 (to remove the - and +1 for the 0) = 71
	TempRedRange = 255.0 / 71.0
	// HumidityBlueRange : 50% humidity for value 123
	HumidityBlueRange = 123.0 / 50.0
)

// WeatherResult : get the elements we need to process from the json response
type WeatherResult struct {
	Temp     float64 `json:"temp"`
	Humidity float64 `json:"humidity"`
	// visibilty contains mainly pollution molecules
	Visibility float64 `json:"visibility"`
	Sky        string
}

func explain(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Please enter the name of the city in the url (ex: Paris)")
}

var mux map[string]func(http.ResponseWriter, *http.Request)

func main() {
	server := http.Server{
		Addr:    ":5000",
		Handler: &myHandler{},
	}
	log.Println("Server started")

	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/"] = explain

	server.ListenAndServe()
}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}
	getWeather(w, r)
}

func (w *WeatherResult) getRed() int {
	// protect the setup from extreme temp
	if w.Temp <= -20 {
		return 0
	} else if w.Temp >= 50 {
		return 255
	}
	// +20 to remove the (-) complexity and format to int
	Red := int((w.Temp + 20) * float64(TempRedRange))
	return Red
}

func (w *WeatherResult) getBlue() int {
	// add 123 if it's raining
	res := int(w.Humidity * HumidityBlueRange)
	if w.Sky == "Rain" {
		res += 123
	}
	return res
}

// 5 stages [0,1,2,3,4] ref: http://www.epa.vic.gov.au/your-environment/air/air-pollution/visibility-reduction
func (w *WeatherResult) getGreen() int {
	switch {
	case w.Visibility >= 45:
		return 255
	case w.Visibility >= 30 && w.Visibility < 45:
		return int(255/4) * 3
	case w.Visibility >= 20 && w.Visibility < 30:
		return int(255/4) * 2
	case w.Visibility >= 10 && w.Visibility < 20:
		return int(255 / 4)
	}
	// case w.Visibility < 10:
	return 0
}

func getMapInterface(weather []byte) map[string]interface{} {
	var dat map[string]interface{}

	if err := json.Unmarshal(weather, &dat); err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	return dat
}

func callOpenWeather(country string) map[string]interface{} {
	response, err := http.Get("http://api.openweathermap.org/data/" + version + "/weather?q=" + country + "&appid=" + key)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		return getMapInterface(contents)
	}
	return nil
}

func getWeather(w http.ResponseWriter, r *http.Request) int {
	// chanel for float values
	t := make(chan float64)
	h := make(chan float64)
	v := make(chan float64)
	// chanel for string main value
	m := make(chan string)

	weather := new(WeatherResult)
	city := strings.Trim(r.URL.Path, "/")
	result := callOpenWeather(city)
	// code is multitype so we are checking the name of the city
	if result["name"] == nil {
		fmt.Fprintf(w, "The name entered is not a known city") // send data to client side
		return -1
	}
	// Get json values
	go utils.SetTemp(result, t)
	go utils.SetHumidity(result, h, m)
	go utils.SetVisibility(result, v)
	// waitting for results
	weather.Temp = <-t
	red := strconv.Itoa(weather.getRed())
	weather.Humidity = <-h
	weather.Sky = <-m
	blue := strconv.Itoa(weather.getBlue())
	weather.Visibility = <-v
	green := strconv.Itoa(weather.getGreen())
	html := "<html><body style=\"background: rgb(" + red + "," + green + "," + blue + ")\"><h1>" + city + "</h1></body></html>"
	fmt.Fprintf(w, html) // send data to client side
	return 0
}
