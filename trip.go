package main

type Trip struct {
	Id         int32   `json:"id"`
	OriginAlt  float64 `json:"origin_alt"`
	OriginLng  float64 `json:"origin_lng"`
	DestinAlt  float64 `json:"destin_alt"`
	DestinLng  float64 `json:"destin_lng"`
	LeaveAfter int32   `json:"leave_after"`
	ArriveBy   int32   `json:"arrive_by"`
	Seats      int8    `json:"seats"`
	DriverUUID string  `json:"driver_uuid"`
}

func NewTrip() *Trip {
	t := new(Trip)
	rows, _ := dbConn.Query("SELECT max(id) FROM trips")
	for rows.Next() {
		rows.Scan(&t.Id)
		break
	}
	t.Id += 1
	return t
}
