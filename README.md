# approaching_109
When buses come near bus stop?
## Sample Output
### Operating Buses
```json
{
  "inOperation": true,
  "message": "List of Buses",
  "busStopName": "大崎広小路",
  "details": [
    {
      "courseName": "渋41",
      "destination": "渋谷駅",
      "stopsBefore": "6",
      "waitAtBusStopMinutes": "11",
      "ETAofBusStop": "18:40",
      "congestion": "空いています",
      "congestionIcon": "congestion-1.png",
      "busNumber": "743"
    },
    {
      "courseName": "渋41",
      "destination": "渋谷駅",
      "stopsBefore": "10",
      "waitAtBusStopMinutes": "20",
      "ETAofBusStop": "18:49",
      "congestion": "空いています",
      "congestionIcon": "congestion-1.png",
      "busNumber": "1028"
    },
    {
      "courseName": "渋41",
      "destination": "渋谷駅",
      "stopsBefore": "始発バス停出発前",
      "waitAtBusStopMinutes": "-380",
      "ETAofBusStop": "12:09",
      "congestion": "",
      "congestionIcon": "",
      "busNumber": "1214"
    }
  ]
}
```
### Last bus has gone
```json
{
  "inOperation": false,
  "message": "Today's last bus is gone. or not provided data.",
  "busStopName": "大崎広小路",
  "details": null
}
```
