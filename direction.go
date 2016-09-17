package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

var mapclient *(maps.Client)

func initClient() *(maps.Client) {
	client, err := maps.NewClient(maps.WithAPIKey(os.Getenv("GOOGLE_API_KEY")))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	mapclient = client
	return client
}

func getTravelTime(origin_alt, origin_lng, destin_alt, destin_lng float64) float64 {
	request := &maps.DirectionsRequest{
		Origin:      fmt.Sprintf("%f,%f", origin_alt, origin_lng),
		Destination: fmt.Sprintf("%f,%f", destin_alt, destin_lng),
	}
	responses, _, err := mapclient.Directions(context.Background(), request)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	return responses[0].Legs[0].Duration.Seconds()
}
