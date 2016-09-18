package main

import (
	"encoding/json"
	"errors"
	"fmt"
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

func GetTripJson(tripid string) string {
	t := new(Trip)
	t.Load(fmt.Sprintf("id = %s", tripid))
	query := fmt.Sprintf("SELECT time, latitude, longitude, action FROM records WHERE trip_id = %s ORDER BY time", tripid)
	rows, _ := dbConn.Query(query)
	stepstr := "\"steps\":["
	for rows.Next() {
		var time int32
		var latitude, longitude float64
		var action string
		rows.Scan(&time, &latitude, &longitude, &action)
		rowstr := fmt.Sprintf("{\"time\":%d,\"latitude\":%f,\"longitude\":%f,\"action\":%s}",
			time, latitude, longitude, action)
		stepstr += rowstr + ","
	}
	jsonmar, _ := json.Marshal(t)
	jsonstr := string(jsonmar)
	jsonstr = jsonstr[:len(jsonstr)-1] + "," + stepstr[:len(stepstr)-1] + "]}"
	return jsonstr
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
