package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Detail struct {
	Congestion  string `json:"congestion"`
	StopsBefore string `json:"stopsBefore"`
	WaitTime    string `json:"waitTime"`
	Estimate    string `json:"estimate"`
	BusNumber   string `json:"busNumber"`
	CourseName  string `json:"courseName"`
	Destination string `json:"destination"`
}

type Bus struct {
	InOperation bool     `json:"inOperation"`
	Details     []Detail `json:"details"`
}

func parseApproaching(url string) Bus {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	bus := Bus{}
	details := []Detail{}
	doc.Find("div#main").Each(func(i int, s *goquery.Selection) {
		if s.Find(".nobusLocationInfo").Length() == 1 {
			// last bus has gone
			bus = Bus{InOperation: false}
			return
		}
	})
	doc.Find("ul#resultList li").Each(func(i int, s *goquery.Selection) {
		// buses will come
		bus = Bus{InOperation: true}
		congestion, _ := s.Find(".congestion-image").Attr("alt")
		stopsBeforeAndWaitText := s.Find(".info").Text()
		stopsAndWaitRex := regexp.MustCompile("[0-9]+個前の停留所を発車【[0-9]+分待ち】")
		stopsBeforeAndWaitTime := stopsAndWaitRex.FindString(stopsBeforeAndWaitText)
		removeTextRex := regexp.MustCompile("(個前の停留所を発車【|分待ち】)")
		stopsAndWait := removeTextRex.Split(stopsBeforeAndWaitTime, -1)
		stopsAndWait = stopsAndWait[:len(stopsAndWait)-1]
		var stopsBefore string
		var waitTime string
		var estimate string
		const layout = "15:04"

		if len(stopsAndWait) > 0 {
			// fmt.Printf("%q %d\n", stopsAndWait, len(stopsAndWait))

			stopsBefore = stopsAndWait[0]
			waitTime = stopsAndWait[1]
			tm := time.Now()
			estimateMin, _ := strconv.Atoi(waitTime)
			tm = tm.Add(time.Duration(estimateMin) * time.Minute)
			estimate = fmt.Sprint(tm.Format(layout))

		} else {
			beforeDepurture := s.Find(".info").Text()
			removeWhiteSpaceRex := regexp.MustCompile(`\s*`)
			beforeDepurture = removeWhiteSpaceRex.ReplaceAllString(beforeDepurture, "")
			stopsBefore = "始発バス停出発前"
			removeHourMinutesRex := regexp.MustCompile("（|）|時|分|に|到着予定")
			arrivalAtStop := removeHourMinutesRex.Split(beforeDepurture, -1)
			arrivalAtStop = arrivalAtStop[1:3]
			arrivalHour, _ := strconv.Atoi(arrivalAtStop[0])
			arrivalMinutes, _ := strconv.Atoi(arrivalAtStop[1])
			loc, _ := time.LoadLocation("Asia/Tokyo")
			tm := time.Now()
			now := time.Date(tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second(), 0, loc)
			arrival := time.Date(tm.Year(), tm.Month(), tm.Day(), arrivalHour, arrivalMinutes, 0, 0, loc)
			duration, _ := time.ParseDuration(arrival.Sub(now).String())
			waitTime = strconv.FormatFloat(duration.Minutes(), 'f', 0, 64)
			estimate = fmt.Sprintf("%02d:%02d", arrivalHour, arrivalMinutes)
		}

		busNumber, _ := s.Find(".locationDataArea #locationData img").Attr("title")
		busNumber = strings.Replace(busNumber, "車両番号:", "", -1)
		courseName := s.Find(".courseName").Text()
		destination := s.Find(".destination-name").Text()

		details = append(details, Detail{
			Congestion:  congestion,
			StopsBefore: stopsBefore,
			WaitTime:    waitTime,
			BusNumber:   busNumber,
			CourseName:  courseName,
			Destination: destination,
			Estimate:    estimate,
		})
	})
	bus.Details = details
	return bus
}

func main() {
	// Buses, _ := json.Marshal(parseApproaching("http://localhost:8888/after_last.html"))
	Buses, _ := json.Marshal(parseApproaching("http://localhost:8888/operating.html"))
	// Buses, _ := json.Marshal(parseApproaching("https://transfer.navitime.biz/tokyubus/smart/location/BusLocationSearchTargetCourse?startId=00240508&poleId=000000001133"))
	fmt.Println(string(Buses))
}
