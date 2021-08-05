# approaching_109
When buses come near bus stop?
## Sample Output
### Operating Buses
```json
{
  "inOperation": true,
  "details": [
    {
      "courseName": "渋４１",
      "destination": "渋谷駅ゆき",
      "stopsBefore": "6",
      "waitAtBusStopTime": "11",
      "estimateArrivalTimeAtStop": "20:41",
      "congestion": "空いています",
      "busNumber": "743"
    },
    {
      "courseName": "渋４１",
      "destination": "渋谷駅ゆき",
      "stopsBefore": "10",
      "waitAtBusStopTime": "20",
      "estimateArrivalTimeAtStop": "20:50",
      "congestion": "空いています",
      "busNumber": "1028"
    },
    {
      "courseName": "渋４１",
      "destination": "渋谷駅ゆき",
      "stopsBefore": "始発バス停出発前",
      "waitAtBusStopTime": "-502",
      "estimateArrivalTimeAtStop": "12:09",
      "congestion": "",
      "busNumber": "1214"
    }
  ]
}
```
### Last bus has gone
```json
{
  "inOperation": false,
  "details": null
}
```
