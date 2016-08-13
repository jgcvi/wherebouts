package main

import (
	"encoding/json"
	"fmt"
	"github.com/GoogleCloudPlatform/golang-samples/getting-started/wherebouts"
	"net/http"
	"strconv"
)

type VendorID struct {
	Uid  string
	Lat  string
	Long string
}

func toggleVendorHandler(w http.ResponseWriter, r *http.Request) *appError {
	decoder := json.NewDecoder(r.Body)
	var v VendorID
	err := decoder.Decode(&v)
	if err != nil {
		fmt.Print("error in reading body")
	}

	uid := v.Uid
	lat, errlat := strconv.ParseFloat(v.Lat, 64)
	long, errlong := strconv.ParseFloat(v.Long, 64)

	if errlat != nil {
		return appErrorf(errlat, "error parsing lat to float")
	}
	if errlong != nil {
		return appErrorf(errlong, "error parsing long to float")
	}

	if err != nil {
		return appErrorf(err, "could not find email: %v", err)
	}

	wherebouts.DB.ToggleAvailability(uid, lat, long)

	w.WriteHeader(200)

	cont := map[string]interface{}{
		"status": "success",
	}

	resp, _ := json.Marshal(cont)
	w.Write(resp)

	return nil
}
