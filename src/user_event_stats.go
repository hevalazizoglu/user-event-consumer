package main

import (
	"net/http"

	"github.com/wcharczuk/go-chart"
)

func ShowUserEventAPIStats(responseWriter http.ResponseWriter, request *http.Request) {
	stats, err := GetResponseTimeStats()

	if err != nil {
		http.Error(responseWriter, "Couldn't fetch response time stats", http.StatusNotFound)
		return
	}

	buckets := make(map[string]float64)
	for _, bucket := range stats {
		buckets[bucket["_id"].(string)] = bucket["count"].(float64)
	}
	sbc := chart.BarChart {
		Height:   512,
		BarWidth: 60,
		XAxis: chart.Style {
			Show: true,
		},
		YAxis: chart.YAxis {
			Style: chart.Style {
					Show: true,
			},
		},
		Bars: []chart.Value {
			{Value: buckets["<1ms"], Label: "<1ms"},
			{Value: buckets["<5ms"], Label: "<5ms"},
			{Value: buckets["<10ms"], Label: "<10ms"},
			{Value: buckets["<20ms"], Label: "<20ms"},
			{Value: buckets["<50ms"], Label: "<50ms"},
			{Value: buckets["<100ms"], Label: "<100ms"},
		},
	}

	responseWriter.Header().Set("Content-Type", "image/png")
	err = sbc.Render(chart.PNG, responseWriter)
	if err != nil {
		http.Error(responseWriter, "Couldn't fetch response time stats", http.StatusNotFound)
		return
	}
}
