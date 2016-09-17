package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Meow!")
}

func RecordCreate(w http.ResponseWriter, r *http.Request) {
	record := NewRecord()
	_ = json.NewDecoder(r.Body).Decode(record)
	record.Save()
}

func TripCreate(w http.ResponseWriter, r *http.Request) {
	t := NewTrip()
	_ = json.NewDecoder(r.Body).Decode(t)
	query := fmt.Sprintf("INSERT INTO trips VALUES(%d, %f, %f, %f, %f, %d, %d, %d, '%s')",
		t.Id, t.OriginAlt, t.OriginLng, t.DestinAlt, t.DestinLng, t.LeaveAfter, t.ArriveBy, t.Seats, t.DriverUUID)

	traveltime := int32(getTravelTime(t.OriginAlt, t.OriginLng, t.DestinAlt, t.DestinLng))

	origin := NewRecord()
	origin.Tripid = t.Id
	origin.Time = t.LeaveAfter
	origin.Altitude = t.OriginAlt
	origin.Longitude = t.OriginLng
	origin.Action = "P " + t.DriverUUID
	origin.Psgcount = 1
	origin.Save()

	destin := NewRecord()
	destin.Tripid = t.Id
	destin.Time = origin.Time + traveltime
	destin.Altitude = t.DestinAlt
	destin.Longitude = t.DestinLng
	destin.Action = "D " + t.DriverUUID
	destin.Psgcount = 1
	destin.Save()

	dbConn.Query(query)
}
