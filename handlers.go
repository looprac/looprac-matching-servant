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
	fmt.Println(t.Id)
	query := fmt.Sprintf("INSERT INTO trips VALUES(%d, %f, %f, %f, %f, %d, %d, %d, '%s')",
		t.Id, t.OriginAlt, t.OriginLng, t.DestinAlt, t.DestinLng, t.LeaveAfter, t.ArriveBy, t.Seats, t.DriverUUID)
	fmt.Println(query)
	dbConn.Query(query)
}
