package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRed(t *testing.T) {
	// test normal conditions
	weather := WeatherResult{20.0, 50, 10, "Clouds"}
	expectedMax := 255
	expectedMin := 0
	actual := weather.getRed()

	if actual < expectedMin && actual > expectedMax {
		t.Errorf("TestGetRed expected to be between 0 and 255 but instead got %d!", actual)
	}

	// test extreme cold temparature
	weather = WeatherResult{-50.0, 50, 10, "Clouds"}
	actual = weather.getRed()

	if actual < expectedMin && actual > expectedMax {
		t.Errorf("TestGetRed expected to be between 0 and 255 but instead got %d!", actual)
	}

	// test extreme hot temparature
	weather = WeatherResult{150.0, 50, 10, "Clouds"}
	actual = weather.getRed()

	if actual < expectedMin && actual > expectedMax {
		t.Errorf("TestGetRed expected to be between 0 and 255 but instead got %d!", actual)
	}
}

func TestGetBlue(t *testing.T) {
	// test max non raining conditions
	weather := WeatherResult{20.0, 50, 10, "Clouds"}
	expected := 123
	actual := weather.getBlue()

	if actual != expected {
		t.Errorf("TestGetBlue expected to be 123 but instead got %d!", actual)
	}

	// test normal rain conditions
	weather = WeatherResult{20.0, 50, 10, "Rain"}
	actual = weather.getBlue()
	expected = 246

	if actual != expected {
		t.Errorf("TestGetBlue expected to be 246 but instead got %d!", actual)
	}

	// test normal condition
	weather = WeatherResult{20.0, 20, 10, "Clouds"}
	actual = weather.getBlue()
	expected = 49

	if actual != expected {
		t.Errorf("TestGetBlue expected to be 49 but instead got %d!", actual)
	}
}

func TestGetGreen(t *testing.T) {
	// test Min conditions
	weather := WeatherResult{20.0, 50, 50, "Clouds"}
	expected := 255
	actual := weather.getGreen()

	if actual != expected {
		t.Errorf("TestGetGreen expected to be 255 but instead got %d!", actual)
	}

	// test Max conditions
	weather = WeatherResult{20.0, 50, 0, "Clouds"}
	expected = 0
	actual = weather.getGreen()

	if actual != expected {
		t.Errorf("TestGetGreen expected to be 0 but instead got %d!", actual)
	}

	// test 10 to 20 conditions
	weather = WeatherResult{20.0, 50, 15, "Clouds"}
	expected = 63
	actual = weather.getGreen()

	if actual != expected {
		t.Errorf("TestGetGreen expected to be 63 but instead got %d!", actual)
	}

	// test case 20 to 30 conditions
	weather = WeatherResult{20.0, 50, 25, "Clouds"}
	expected = 126
	actual = weather.getGreen()

	if actual != expected {
		t.Errorf("TestGetGreen expected to be 126 but instead got %d!", actual)
	}

	// test 30 to 40 conditions
	weather = WeatherResult{20.0, 50, 35, "Clouds"}
	expected = 189
	actual = weather.getGreen()

	if actual != expected {
		t.Errorf("TestGetGreen expected to be 189 but instead got %d!", actual)
	}
}

func TestGetWeather(t *testing.T) {
	// test wrong name city
	req, err := http.NewRequest("GET", "/WrongNameCity", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	actual := getWeather(rr, req)
	expected := -1

	if actual != expected {
		t.Errorf("TestGetWeather expected to be -1 but instead got %d!", actual)
	}

	// test Paris name city
	req, err = http.NewRequest("GET", "/Paris", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()

	actual = getWeather(rr, req)
	expected = 0

	if actual != expected {
		t.Errorf("TestGetWeather expected to be 0 but instead got %d!", actual)
	}
}
