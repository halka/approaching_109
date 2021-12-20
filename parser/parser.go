package parser

import (
	"approaching_109/model"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/unicode/norm"
)

func SetParams(startId string, poleId string) model.Params {
	p := model.Params{}
	p.StartId = startId
	p.PoleId = poleId

	return p
}

func ParseApproaching(p model.Params) model.Bus {
	url := "https://transfer.navitime.biz/tokyubus/smart/location/BusLocationSearchTargetCourse?startId=" + p.StartId + "&poleId=" + p.PoleId
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
	bus := model.Bus{}
	var busStopName string
	var details []model.Detail
	doc.Find("div#main").Each(func(i int, s *goquery.Selection) {
		busStopName = s.Find("#allSummary div .big-font").Text()
		if s.Find(".nobusLocationInfo").Length() == 1 {
			// last bus has gone
			bus = model.Bus{
				InOperation: false,
				Message:     "Today's last bus is gone. or not provided data.",
				BusStopName: busStopName,
			}
			return
		}
	})
	doc.Find("ul#resultList li").Each(func(i int, s *goquery.Selection) {
		// buses will come
		bus = model.Bus{
			InOperation: true,
			Message:     "List of Buses",
			BusStopName: busStopName,
		}
		congestion, _ := s.Find(".congestion-image").Attr("alt")
		congestionIcon, _ := s.Find(".congestion-image").Attr("src")
		congestionIcon = strings.Replace(string(congestionIcon), "/blt-storage/pc/img/tokyubus/location/", "", -1)
		congestionLevelRex := regexp.MustCompile("[0-9]+")
		congestionLevel := congestionLevelRex.FindString(congestionIcon)
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
			stopsBefore = stopsAndWait[0]
			waitTime = stopsAndWait[1]
			tm := time.Now()
			estimateMin, _ := strconv.Atoi(waitTime)
			tm = tm.Add(time.Duration(estimateMin) * time.Minute)
			estimate = fmt.Sprint(tm.Format(layout))

		} else {
			beforeDeparture := s.Find(".info").Text()
			removeWhiteSpaceRex := regexp.MustCompile(`\s*`)
			beforeDeparture = removeWhiteSpaceRex.ReplaceAllString(beforeDeparture, "")
			stopsBefore = "始発バス停出発前"
			removeHourMinutesRex := regexp.MustCompile("（|）|時|分|に|到着予定")
			arrivalAtStop := removeHourMinutesRex.Split(beforeDeparture, -1)
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
		courseName := string(norm.NFKC.Bytes([]byte(s.Find(".courseName").Text())))
		destination := strings.Replace(s.Find(".destination-name").Text(), "ゆき", "", -1)

		details = append(details, model.Detail{
			Congestion:         congestion,
			CongestionIcon:     congestionIcon,
			CongrestionLevelse: congestionLevel,
			StopsBefore:        stopsBefore,
			WaitAtBusStop:      waitTime,
			BusNumber:          busNumber,
			CourseName:         courseName,
			Destination:        destination,
			ETAofBusStop:       estimate,
		})
	})
	bus.Details = details
	return bus
}
