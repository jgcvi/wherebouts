package main

import (
	"encoding/json"
	"fmt"
	"github.com/GoogleCloudPlatform/golang-samples/getting-started/wherebouts"
	//"github.com/satori/go.uuid"
	"net/http"
)

type registeree struct {
	Name     string
	Password string
	Email    string
}

type LoginAttempt struct {
	Email    string
	Password string
}

func loginHandler(w http.ResponseWriter, r *http.Request) *appError {
	decoder := json.NewDecoder(r.Body)
	var l LoginAttempt
	err := decoder.Decode(&l)
	if err != nil {

	}

	pass := l.Password
	email := l.Email

	vendor, err := wherebouts.DB.GetVendorByEmail(email)
	if err != nil {
		return appErrorf(err, "could not find email: %v", err)
	}

	if !vendor.ValidatePass(pass) {
		return appErrorf(nil, "password or email wrong")
	}

	w.WriteHeader(200)

	resp, _ := json.Marshal(vendor)
	fmt.Print(string(resp))
	w.Write(resp)

	return nil
}

func logoutHandler(w http.ResponseWriter, r *http.Request) *appError {
	return nil
}

func registerHandler(w http.ResponseWriter, r *http.Request) *appError {
	decoder := json.NewDecoder(r.Body)
	var t registeree
	err := decoder.Decode(&t)
	if err != nil {
		return appErrorf(nil, "could not decode user")
	}
	name := t.Name
	pass := t.Password
	email := t.Email

	if name != "" && pass != "" && email != "" {
		v := wherebouts.CreateVendor(name, email, pass)
		wherebouts.DB.AddVendor(v)

		w.WriteHeader(200)

		cont := map[string]interface{}{
			"status": "success",
		}
		resp, _ := json.Marshal(cont)

		w.Write(resp)

		return nil

	} else if pass == "" {
		return appErrorf(nil, "pass was empty")
	} else if name == "" {
		return appErrorf(nil, "name empty")
	} else if email == "" {
		return appErrorf(nil, "email empty")
	}

	return appErrorf(nil, "could not create entry")
}
