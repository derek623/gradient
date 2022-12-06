package customer_service

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"reflect"
	"sort"
	"strconv"

	"git.codesubmit.io/sfox/party-invite-ruiegv/pkg/greatCircle"
	"git.codesubmit.io/sfox/party-invite-ruiegv/pkg/util"
)

var OfficeLocation greatCircle.Point

// Set the office location
func SetOfficeLocation(officeLongitude float64, officeLatitude float64) error {
	OfficeLocation.Longitude = greatCircle.DegreeToRadian(officeLongitude)
	OfficeLocation.Latitude = greatCircle.DegreeToRadian(officeLatitude)
	log.Println("Set office location to ", OfficeLocation)
	if !OfficeLocation.Valid() {
		return errors.New("Invalid office longitude or latitude")
	}

	return nil
}

// Customer struct to store customer information
type Customer struct {
	Latitude  string
	User_id   int
	Name      string
	Longitude string
	Location  greatCircle.Point
}

// Implement MarshalJSON for Customer to only print user id and name
func (c *Customer) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		User_id int
		Name    string
	}{
		c.User_id, c.Name,
	})
}

// helper function to check if the unmarshaled JSON has all required keys
func hasRequiredKey(m map[string]interface{}) bool {
	_, longitudeFound := m["longitude"]
	_, latitudeFound := m["latitude"]
	_, userIdFound := m["user_id"]
	_, nameFound := m["name"]

	return longitudeFound && latitudeFound && userIdFound && nameFound
}

// Implement the UnmarshalJSON function so as to perform proper checking on the JSON and
// convert the longitude and latitude into radian
func (c *Customer) UnmarshalJSON(b []byte) error {
	var tmpCustomer map[string]interface{}
	if err := json.Unmarshal(b, &tmpCustomer); err != nil {
		return err
	}
	if !hasRequiredKey(tmpCustomer) {
		return errors.New("JSON missing required keys, please check")
	}
	//Now try to convert the values into the appropriate type
	for key, value := range tmpCustomer {
		switch key {
		case "longitude":
			if reflect.TypeOf(value).Kind() != reflect.String {
				return errors.New("Cannot convert longitude as value is not of type string")
			}
			c.Longitude = value.(string)
		case "latitude":
			if reflect.TypeOf(value).Kind() != reflect.String {
				return errors.New("Cannot convert latitude as value is not of type string")
			}
			c.Latitude = value.(string)
		case "user_id":
			if reflect.TypeOf(value).Kind() != reflect.Float64 {
				return errors.New("Cannot convert user_id as value is not of type float64")
			}
			c.User_id = int(value.(float64))
		case "name":
			if reflect.TypeOf(value).Kind() != reflect.String {
				return errors.New("Cannot convert name as value is not of type string")
			}
			c.Name = value.(string)

		}
	}

	longtitude, err := strconv.ParseFloat(c.Longitude, 64)
	if nil != err {
		return err
	}

	latitude, err := strconv.ParseFloat(c.Latitude, 64)
	if nil != err {
		return err
	}

	p := greatCircle.MakePoint(greatCircle.DegreeToRadian(longtitude), greatCircle.DegreeToRadian(latitude))
	if !p.Valid() {
		return errors.New(" Invalid longitude or latitude")
	}

	c.Location = p

	return nil
}

// Test if we should invite the customer
func (customer Customer) shouldInviteCustomer(distance float64) (bool, error) {
	if distance < 0.0 {
		return false, errors.New("Distance must be > 0")
	}

	return util.SmallerOrEqual(distance, 100.0), nil
}

// Convert byte array into a customer map
func convertToCustomers(filebyte []byte) (map[int]Customer, error) {

	//Split by the newline character
	customers := bytes.Split(filebyte, []byte("\n"))
	customerMap := make(map[int]Customer, len(customers))
	for _, customer := range customers {
		if !json.Valid(customer) {
			err := "Invalid JSON: " + string(customer)
			return nil, errors.New(err)
		}

		var c Customer
		err := json.Unmarshal([]byte(customer), &c)
		if nil != err {
			err := "Cannot unmarshal customer " + string(customer) + " : " + err.Error()
			return nil, errors.New(err)
		}
		_, found := customerMap[c.User_id]
		if found {
			err := "Customer id overlap: " + strconv.Itoa(c.User_id)
			return nil, errors.New(err)
		}
		customerMap[c.User_id] = c

	}
	return customerMap, nil
}

// Entry point of the customer service
func GetCustomers(w http.ResponseWriter, r *http.Request) error {

	if http.MethodPut != r.Method {
		return errors.New("HTTP request is not a PUT request")
	}
	fileBytes, err := util.GetFileInBytes(r)
	if nil != err {
		return err
	}

	customers, err := convertToCustomers(fileBytes)
	if nil != err {
		return err
	}
	log.Println(customers)

	//invite the appropriate customers
	var resultCustomerSlice []int
	for _, customer := range customers {
		result, err := customer.shouldInviteCustomer(greatCircle.Distance(OfficeLocation, customer.Location, greatCircle.Radius))
		if err != nil {
			return err
		}
		if result {
			resultCustomerSlice = append(resultCustomerSlice, customer.User_id)
		}
	}

	//sort the result slice
	sort.Ints(resultCustomerSlice)
	var sortedResultCustomerSlice []Customer
	for _, key := range resultCustomerSlice {
		sortedResultCustomerSlice = append(sortedResultCustomerSlice, customers[key])
	}

	//return results in JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp, err := json.Marshal(&sortedResultCustomerSlice)
	if nil != err {
		return err
	}

	w.Write(resp)

	return nil
}
