// Package api provides structure of the api server
package api

import (
	"log"
	"net/http"

	"git.codesubmit.io/sfox/party-invite-ruiegv/pkg/customer_service"
	"git.codesubmit.io/sfox/party-invite-ruiegv/pkg/util"
)

// This is the version 1 struct
type ApiV1 struct {
}

// Register the provided handle
func (api *ApiV1) registerHandle(patterns string, f func(w http.ResponseWriter, r *http.Request)) {
	http.HandleFunc(patterns, f)
}

// Return the existing version
func (api *ApiV1) getVersion() string {
	return "v1"
}

// Return an apiV1 struct with everything initialized (e.g, office location initialized and proper handle registered)
func GetApiV1(officeLongitude float64, officeLatitude float64) (*ApiV1, error) {
	if err := customer_service.SetOfficeLocation(officeLongitude, officeLatitude); err != nil {
		return nil, err
	}
	api := &ApiV1{}
	pattern := "/" + api.getVersion() + "/customer"
	api.registerHandle(pattern, util.ErrorHandler(customer_service.GetCustomers))
	return api, nil
}

// Start up the server and listen to the provided addr
func (api *ApiV1) StartServer(addr string) error {
	log.Println("Starting server...")
	if err := http.ListenAndServe(addr, nil); err != nil {
		return err
	}
	return nil
}
