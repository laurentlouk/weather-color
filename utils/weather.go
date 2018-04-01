package utils

// SetTemp : get temp in T(°C)
func SetTemp(res map[string]interface{}, c chan float64) {
	main := res["main"]
	temp := main.(map[string]interface{})["temp"]
	// convert to °C -> T(°C) = T(K) - 273.15
	c <- temp.(float64) - 273.15
}

// SetHumidity : get humity that can't exceed 50 & weather condition (Clouds, Rain, ...)
func SetHumidity(res map[string]interface{}, c chan float64, m chan string) {
	main := res["weather"]
	index := main.([]interface{})[0]
	// get sky if it's raining
	sky := index.(map[string]interface{})["main"].(string)
	humidity := res["main"].(map[string]interface{})["humidity"].(float64)
	// humidity can't exceed 50%
	if humidity > 50 {
		humidity = 50
	}
	c <- humidity
	m <- sky
}

// SetVisibility : get visibility for pollution information
func SetVisibility(res map[string]interface{}, c chan float64) {
	// m to km
	c <- res["visibility"].(float64) / 1000
}
