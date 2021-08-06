package model

type Detail struct {
	CourseName                   string `json:"courseName"`
	Destination                  string `json:"destination"`
	StopsBefore                  string `json:"stopsBefore"`
	WaitAtBusStop                string `json:"waitAtBusStopMinutes"`
	ETAofBusStop 				 string `json:"ETAofBusStop"`
	Congestion                   string `json:"congestion"`
	CongestionIcon               string `json:"congestionIcon"`
	BusNumber                    string `json:"busNumber"`
}

type Bus struct {
	InOperation bool     `json:"inOperation"`
	Message		string	 `json:"message"`
	BusStopName	string	 `json:"busStopName"`
	Details     []Detail `json:"details"`
}

type Params struct {
	StartId, PoleId string
}
