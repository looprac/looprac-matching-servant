package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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
		t.Id, t.OriginLat, t.OriginLng, t.DestinLat, t.DestinLng, t.LeaveAfter, t.ArriveBy, t.Seats, t.DriverUUID)
	dbConn.QueryRow(query)

	traveltime := int32(getTravelTime(t.OriginLat, t.OriginLng, t.DestinLat, t.DestinLng))

	origin := NewRecord()
	origin.Tripid = t.Id
	origin.Time = t.LeaveAfter
	origin.Latitude = t.OriginLat
	origin.Longitude = t.OriginLng
	origin.Action = "P " + t.DriverUUID
	origin.Psgcount = 1
	origin.Save()

	destin := NewRecord()
	destin.Tripid = t.Id
	destin.Time = origin.Time + traveltime
	destin.Latitude = t.DestinLat
	destin.Longitude = t.DestinLng
	destin.Action = "D " + t.DriverUUID
	destin.Psgcount = 1
	destin.Save()

}

func TripGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tripid := vars["tripid"]
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
	fmt.Fprintln(w, jsonstr)
}
