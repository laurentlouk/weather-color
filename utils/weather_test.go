package utils

import (
	"encoding/json"
	"fmt"
	"testing"
)

// json to map:intefrace for tests
func getMapInterface(weather []byte) map[string]interface{} {
	var dat map[string]interface{}

	if err := json.Unmarshal(weather, &dat); err != nil {
		fmt.Printf("%s", err)
	}
	return dat
}

// EPSILON :float64 equality esitimation
var EPSILON = 0.00000001

func floatEquals(a, b float64) bool {
	if (a-b) < EPSILON && (b-a) < EPSILON {
		return true
	}
	return false
}

func TestSetTemp(t *testing.T) {
	// test temp from T(K) to T(°C)
	temp := make(chan float64)
	bjson := []byte(`{"weather":[{"main":"Rain"}],"main":{"temp":290.30,"humidity":72},"visibility":10000,"name":"Mykonos"}`)
	result := getMapInterface(bjson)
	expected := float64(17.150000)
	go SetTemp(result, temp)
	actual := <-temp

	res := floatEquals(expected, actual)

	if res != true {
		t.Errorf("TestSetTemp expected to be %f but instead got %f!", expected, actual)
	}

	// test minus temp from T(K) to T(°C)
	temp = make(chan float64)
	bjson = []byte(`{"weather":[{"main":"Rain"}],"main":{"temp":90.30,"humidity":72},"visibility":10000,"name":"Mykonos"}`)
	result = getMapInterface(bjson)
	expected = float64(-182.850000)
	go SetTemp(result, temp)
	actual = <-temp

	res = floatEquals(expected, actual)

	if res != true {
		t.Errorf("TestSetTemp expected to be %f but instead got %f!", expected, actual)
	}
}

func TestSetHumidity(t *testing.T) {
	// test with normal humidity + Rain value
	hum := make(chan float64)
	m := make(chan string)
	bjson := []byte(`{"weather":[{"main":"Rain"}],"main":{"temp":90.30,"humidity":32},"visibility":10000,"name":"Mykonos"}`)
	result := getMapInterface(bjson)
	expectedM := "Rain"
	expectedHum := float64(32.0)
	go SetHumidity(result, hum, m)
	actualHum := <-hum
	actualM := <-m

	if actualM != expectedM {
		t.Errorf("TestSetHumidity expected to be %s but instead got %s!", expectedM, actualM)
	}

	res := floatEquals(actualHum, expectedHum)

	if res != true {
		t.Errorf("TestSetHumidity expected to be %f but instead got %f!", expectedHum, actualHum)
	}

	// test with more than 50% humidity
	hum = make(chan float64)
	m = make(chan string)
	bjson = []byte(`{"weather":[{"main":"Rain"}],"main":{"temp":90.30,"humidity":82},"visibility":10000,"name":"Mykonos"}`)
	result = getMapInterface(bjson)
	expectedM = "Rain"
	expectedHum = float64(50.0)
	go SetHumidity(result, hum, m)
	actualHum = <-hum
	actualM = <-m

	res = floatEquals(actualHum, expectedHum)

	if res != true {
		t.Errorf("TestSetHumidity expected to be %f but instead got %f!", expectedHum, actualHum)
	}
}

func TestSetVisibility(t *testing.T) {
	// test visibilty made of PM2.5 and PM10 particles to who contains pollution
	vis := make(chan float64)
	bjson := []byte(`{"weather":[{"main":"Rain"}],"main":{"temp":90.30,"humidity":82},"visibility":10000,"name":"Mykonos"}`)
	expected := float64(10.0)
	result := getMapInterface(bjson)
	go SetVisibility(result, vis)
	actual := <-vis

	res := floatEquals(actual, expected)

	if res != true {
		t.Errorf("TestSetHumidity expected to be %f but instead got %f!", expected, actual)
	}

}
