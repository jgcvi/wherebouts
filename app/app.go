package main

import (
	"encoding/json"
	//	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	//	"path"
	//	"golang.org/x/net/context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	//	"github.com/satori/go.uuid"
	//	"google.golang.org/cloud/storage"
	"github.com/GoogleCloudPlatform/golang-samples/getting-started/wherebouts"
)

func main() {
	r := mux.NewRouter()

	r.Handle("/", http.RedirectHandler("/vendors", http.StatusFound))

	r.Methods("GET").Path("/vendors").
		Handler(appHandler(listHandler))
	//	r.Methods("GET").Path("/vendors/myFavorite").
	//		Handler(appHandler(listFavoriteHandler))

	r.Methods("GET").Path("/nearbyVendors").
		Handler(appHandler(nearbyHandler))

	r.Methods("GET").Path("/_ah/health").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})

	r.Methods("POST").Path("/login").
		Handler(appHandler(loginHandler))

	r.Methods("POST").Path("/logout").
		Handler(appHandler(logoutHandler))

	r.Methods("POST").Path("/register").
		Handler(appHandler(registerHandler))

	r.Methods("POST").Path("/toggle").
		Handler(appHandler(toggleVendorHandler))

	http.Handle("/", handlers.CombinedLoggingHandler(os.Stderr, r))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func listHandler(w http.ResponseWriter, r *http.Request) *appError {
	vendors, err := wherebouts.DB.ListVendors()

	// TODO handle error case here
	if err != nil {
		return appErrorf(err, "could not list vendors: %v", err)
	}

	sliced, _ := json.Marshal(vendors)

	w.Write(sliced)

	return nil
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

type appError struct {
	Error   error
	Message string
	Code    int
}

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		log.Printf("Handler error: status code: %d, message: %s, underlying err: %#v",
			e.Code, e.Message, e.Error)

		http.Error(w, e.Message, e.Code)
	}
}

func appErrorf(err error, format string, v ...interface{}) *appError {
	return &appError{
		Error:   err,
		Message: fmt.Sprintf(format, v...),
		Code:    500,
	}
}
