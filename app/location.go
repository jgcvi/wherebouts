package main

import (
	"encoding/json"
	"fmt"
	"github.com/GoogleCloudPlatform/golang-samples/getting-started/wherebouts"
	"net/http"
	"strconv"
)

func nearbyHandler(w http.ResponseWriter, r *http.Request) *appError {
	r.ParseForm()
	lat := r.FormValue("latitude")
	long := r.FormValue("longitude")
	dist := r.FormValue("distance")

	parLat, _ := strconv.ParseFloat(lat, 64)
	parLong, _ := strconv.ParseFloat(long, 64)
	parDist, _ := strconv.ParseFloat(dist, 64)

	fmt.Printf("In location: %v, %v, %v\n", parLat, parLong, parDist)

	vendors, err := wherebouts.DB.GetNearbyVendors(parLat, parLong, parDist)

	if err != nil {
		appErrorf(err, "could not find vendors: %v", err)
	}
	w.WriteHeader(200)

	resp, _ := json.Marshal(vendors)
	w.Write(resp)

	return nil
}
