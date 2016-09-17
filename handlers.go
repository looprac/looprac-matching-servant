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
