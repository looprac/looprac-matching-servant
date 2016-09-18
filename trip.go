package main

import (
	"errors"
	"log"
)

type Trip struct {
	Id         int32   `json:"id"`
	OriginLat  float64 `json:"origin_lat"`
	OriginLng  float64 `json:"origin_lng"`
	DestinLat  float64 `json:"destin_lat"`
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

func (t *Trip) Load(condition string) error {
	rows, err := dbConn.Query("SELECT * FROM trips WHERE " + condition)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err := rows.Scan(&t.Id, &t.OriginLat, &t.OriginLng, &t.DestinLat, &t.DestinLng,
			&t.LeaveAfter, &t.ArriveBy, &t.Seats, &t.DriverUUID)
		return err
	}
	return errors.New("Not Found")
}
