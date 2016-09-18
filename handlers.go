package main

import (
	"encoding/json"
	"fmt"
	"log"
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

	traveltime := int32(getTravelTime(t.OriginLat, t.OriginLng, t.DestinLat, t.DestinLng))
	if traveltime < 0 {
		log.Println("Invalid Localtion")
		return
	}
	log.Println("Travel time: ", traveltime)

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

	dbConn.QueryRow(query)
}

func TripGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tripid := vars["tripid"]
	jsonstr := GetTripJson(tripid)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, jsonstr)
}

func matchTrips(lat, lng float64, leave_after, arrive_by int32) []int32 {
	log.Println(lat, lng, leave_after, arrive_by)
	geoeps := 0.1

	query := "SELECT trip_id FROM records"
	query += fmt.Sprintf(" WHERE ABS(latitude - %f) < %f", lat, geoeps)
	query += fmt.Sprintf(" AND ABS(longitude - %f) < %f", lng, geoeps)
	//query += fmt.Sprintf(" AND time <= %d", leave_after)
	//query += fmt.Sprintf(" AND time >= %d", arrive_by)
	query += " ORDER BY trip_id"
	rows, _ := dbConn.Query(query)
	ids := []int32{}
	for rows.Next() {
		var id int32
		rows.Scan(&id)
		ids = append(ids, id)
	}
	return ids
}

func unionSet(id1, id2 []int32) []int32 {
	ids := []int32{}
	l1 := 0
	l2 := 0
	r1 := len(id1)
	r2 := len(id2)
	for l1 < r1 || l2 < r2 {
		if l1 < r1 && l2 < r2 && id1[l1] == id2[l2] {
			ids = append(ids, id2[l2])
			l1 += 1
			l2 += 1
		} else if l2 == r2 || (l1 < r1 && id1[l1] < id2[l2]) {
			l1 += 1
		} else {
			l2 += 1
		}
	}
	return ids
}

func TripSearch(w http.ResponseWriter, r *http.Request) {
	log.Println("Search")
	t := new(Trip)
	_ = json.NewDecoder(r.Body).Decode(t)

	id1 := matchTrips(t.OriginLat, t.OriginLng, t.LeaveAfter, t.ArriveBy)
	id2 := matchTrips(t.DestinLat, t.DestinLng, t.LeaveAfter, t.ArriveBy)
	ids := unionSet(id1, id2)
	log.Println(id1, id2)

	trip_info := "["
	for _, id := range ids {
		trip_info += GetTripJson(fmt.Sprint(id)) + ","
	}
	if len(trip_info) > 1 {
		trip_info = trip_info[:len(trip_info)-1]
	}
	trip_info += "]"
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, trip_info)
}
