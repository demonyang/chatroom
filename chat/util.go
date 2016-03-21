package chat

import (
    "fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type Weather struct {
    Desc   string
	Status int
	Data   Datainfo
}

type Datainfo struct {
    Wendu   string
	Ganmao  string
	Forecast []Forcastinfo
	Yes     Yesinfo
	Aqi     string
	City    string
}

type Yesinfo struct {
    Fl   string
	Fx   string
	High string
	Type string
	Low  string
	Date string
}

type Forcastinfo struct {
    Fengxiang  string
	Fengli     string
	High       string
	Type       string
	Low        string
	Date       string
}

const (
	weatherurl = "http://wthrcdn.etouch.cn/weather_mini?city="
)

func  (s *Server) GetWeather(city string) (string, error) {
	var weatherinfo Weather
	strurl := weatherurl + city
    resp, err := http.Get(strurl,)
	if err != nil {
		return "", err
	}
    defer resp.Body.Close()
	input, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	    return "", err
	}
	if err := json.Unmarshal(input, &weatherinfo); err != nil {
	    return "", err
	}
	// weatherinfo not null
	if weatherinfo.Desc == "OK" {
		weastr := fmt.Sprintf("城市:%v,日期:%v,天气:%v,温度:%v-%v",weatherinfo.Data.City, weatherinfo.Data.Forecast[0].Date, weatherinfo.Data.Forecast[0].Type, weatherinfo.Data.Forecast[0].Low, weatherinfo.Data.Forecast[0].High)
	    return weastr, nil
	}
	return "", fmt.Errorf("Get weather failed")
}
