package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/satori/go.uuid"
)

type Record struct {
	Tripid    string  `json:"trip_id"`
	Time      int32   `json:"time"`
	Altitude  float64 `json:"altitude"`
	Longitude float64 `json:"longitude"`
	Action    string  `json:"action"`
	Psgcount  int8    `json:"psg_count"`
	uuid      string  `json:"uuid"`
}

func NewRecord() *Record {
	r := new(Record)
	r.uuid = fmt.Sprintf("%s", uuid.NewV4())
	return r
}

func (r *Record) Print() {
	fmt.Println(r.Tripid, r.Time, r.Altitude, r.Longitude, r.Action, r.Psgcount)
}

func (r *Record) Load(condition string) error {
	rows, err := dbConn.Query("SELECT * FROM records WHERE " + condition)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err := rows.Scan(&r.Tripid, &r.Time, &r.Altitude, &r.Longitude, &r.Action, &r.Psgcount, &r.uuid)
		return err
	}
	return errors.New("Not Found")
}

func (r *Record) Save() error {
	existed := new(Record)
	err := existed.Load(fmt.Sprintf("uuid = '%s'", r.uuid))
	if err == nil {
		query := fmt.Sprintf("UPDATE records SET (Time, Altitude, Longitude, Action, psg_count) = (%d, %f, %f, '%s', %d) WHERE uuid = '%s'",
			r.Time, r.Altitude, r.Longitude, r.Action, r.Psgcount, existed.uuid)
		fmt.Println(query)
		dbConn.QueryRow(query)
	} else {
		query := fmt.Sprintf("INSERT INTO records VALUES('%s', %d, %f, %f, '%s', %d, '%s')",
			r.Tripid, r.Time, r.Altitude, r.Longitude, r.Action, r.Psgcount, r.uuid)
		fmt.Println(query)
		dbConn.QueryRow(query)
	}
	return nil
}
