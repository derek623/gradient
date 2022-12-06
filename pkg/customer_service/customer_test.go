package customer_service

import (
	"errors"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"git.codesubmit.io/sfox/party-invite-ruiegv/pkg/greatCircle"
	"git.codesubmit.io/sfox/party-invite-ruiegv/pkg/util"
)

type shouldInviteCustomerTest struct {
	distance float64
	expected bool
	err      error
}

var shouldInviteCustomerTests []shouldInviteCustomerTest = []shouldInviteCustomerTest{
	shouldInviteCustomerTest{0, true, nil},
	shouldInviteCustomerTest{0.1, true, nil},
	shouldInviteCustomerTest{10, true, nil},
	shouldInviteCustomerTest{50, true, nil},
	shouldInviteCustomerTest{99.9, true, nil},
	shouldInviteCustomerTest{100, true, nil},
	shouldInviteCustomerTest{100.1, false, nil},
	shouldInviteCustomerTest{-1, false, errors.New("Distance must be > 0")},
	shouldInviteCustomerTest{-49.999, false, errors.New("Distance must be > 0")},
}

func TestShouldInviteCustomer(t *testing.T) {
	var c Customer
	for _, test := range shouldInviteCustomerTests {
		if b, err := c.shouldInviteCustomer(test.distance); b != test.expected || (err != nil && err.Error() != test.err.Error()) {
			t.Errorf("Output %v not equal to expected %v", b, test.expected)
		}
	}

}

type marshalJSONTest struct {
	userid   int
	name     string
	expected string
}

var marshalJSONTests []marshalJSONTest = []marshalJSONTest{
	marshalJSONTest{0, "", "{\"User_id\":0,\"Name\":\"\"}"},
	marshalJSONTest{1, "", "{\"User_id\":1,\"Name\":\"\"}"},
	marshalJSONTest{0, "John", "{\"User_id\":0,\"Name\":\"John\"}"},
	marshalJSONTest{1, "John", "{\"User_id\":1,\"Name\":\"John\"}"},
	marshalJSONTest{600, "!@#$%^*()", "{\"User_id\":600,\"Name\":\"!@#$%^*()\"}"},
}

func TestMarshalJSON(t *testing.T) {
	var c Customer
	for _, test := range marshalJSONTests {
		c.User_id = test.userid
		c.Name = test.name
		if b, _ := c.MarshalJSON(); string(b) != test.expected {
			t.Errorf("Output %v not equal to expected %v", b, test.expected)
		}
	}

}

type hasRequiredKeyTest struct {
	m        map[string]interface{}
	expected bool
}

var hasRequiredKeyTests []hasRequiredKeyTest = []hasRequiredKeyTest{
	{map[string]interface{}{}, false},
	{map[string]interface{}{"latitude": ""}, false},
	{map[string]interface{}{"latitude": "", "longitude": ""}, false},
	{map[string]interface{}{"latitude": "", "longitude": "", "user_id": ""}, false},
	{map[string]interface{}{"latitude": "", "longitude": "", "user_id": "", "name": ""}, true},
}

func TestHasRequiredKeys(t *testing.T) {
	for _, test := range hasRequiredKeyTests {
		if result := hasRequiredKey(test.m); result != test.expected {
			t.Errorf("Output %v not equal to expected %v", result, test.expected)
		}
	}
}

type unmarshalJSONTest struct {
	input     string
	c         Customer
	errString string
}

var unmarshalJSONTests []unmarshalJSONTest = []unmarshalJSONTest{
	unmarshalJSONTest{"{\"latitude\": \"52.986375\", \"user_id\": 12, \"name\": \"Christina McArdle\", \"longitude\": \"-6.043701\"}",
		Customer{"52.986375", 12, "Christina McArdle", "-6.043701", greatCircle.MakePoint(greatCircle.DegreeToRadian(-6.043701), greatCircle.DegreeToRadian(52.986375))},
		""},
	unmarshalJSONTest{"{\"latitude\": \"u\", \"user_id\": 12, \"name\": \"Christina McArdle\", \"longitude\": \"-6.043701\"}",
		Customer{"52.986375", 12, "Christina McArdle", "-6.043701", greatCircle.MakePoint(greatCircle.DegreeToRadian(-6.043701), greatCircle.DegreeToRadian(52.986375))},
		strconv.ErrSyntax.Error()},
	unmarshalJSONTest{"{\"latitude\": \"52.986375\", \"user_id\": 12, \"name\": \"Christina McArdle\", \"longitude\": \"iii\"}",
		Customer{"52.986375", 12, "Christina McArdle", "-6.043701", greatCircle.MakePoint(greatCircle.DegreeToRadian(-6.043701), greatCircle.DegreeToRadian(52.986375))},
		strconv.ErrSyntax.Error()},
	unmarshalJSONTest{"{\"latitude\": \"-91\", \"user_id\": 12, \"name\": \"Christina McArdle\", \"longitude\": \"-6.043701\"}",
		Customer{"52.986375", 12, "Christina McArdle", "-6.043701", greatCircle.MakePoint(greatCircle.DegreeToRadian(-6.043701), greatCircle.DegreeToRadian(52.986375))},
		"Invalid longitude or latitude"},
	unmarshalJSONTest{"{\"latitude\": \"-90\", \"user_id\": 12, \"name\": \"Christina McArdle\", \"longitude\": \"181\"}",
		Customer{"52.986375", 12, "Christina McArdle", "-6.043701", greatCircle.MakePoint(greatCircle.DegreeToRadian(-6.043701), greatCircle.DegreeToRadian(52.986375))},
		"Invalid longitude or latitude"},
	unmarshalJSONTest{"{\"latitude\": 52.986375, \"user_id\": 12, \"name\": \"Christina McArdle\", \"longitude\": \"-6.043701\"}",
		Customer{"52.986375", 12, "Christina McArdle", "-6.043701", greatCircle.MakePoint(greatCircle.DegreeToRadian(-6.043701), greatCircle.DegreeToRadian(52.986375))},
		"Cannot convert latitude"},
	unmarshalJSONTest{"{\"latitude\": \"52.986375\", \"user_id\": 12, \"name\": \"Christina McArdle\", \"longitude\": -6.043701}",
		Customer{"52.986375", 12, "Christina McArdle", "-6.043701", greatCircle.MakePoint(greatCircle.DegreeToRadian(-6.043701), greatCircle.DegreeToRadian(52.986375))},
		"Cannot convert longitude"},
	unmarshalJSONTest{"{\"latitude\": \"52.986375\", \"user_id\": \"12\", \"name\": \"Christina McArdle\", \"longitude\": \"-6.043701\"}",
		Customer{"52.986375", 12, "Christina McArdle", "-6.043701", greatCircle.MakePoint(greatCircle.DegreeToRadian(-6.043701), greatCircle.DegreeToRadian(52.986375))},
		"Cannot convert user_id"},
	unmarshalJSONTest{"{\"latitude\": \"52.986375\", \"user_id\": 12, \"name\": 34534, \"longitude\": \"-6.043701\"}",
		Customer{"52.986375", 12, "Christina McArdle", "-6.043701", greatCircle.MakePoint(greatCircle.DegreeToRadian(-6.043701), greatCircle.DegreeToRadian(52.986375))},
		"Cannot convert name"},
}

func TestUnmarshalJSON(t *testing.T) {
	var c Customer
	for _, test := range unmarshalJSONTests {
		err := c.UnmarshalJSON([]byte(test.input))
		if err != nil {
			if !strings.Contains(err.Error(), test.errString) {
				t.Errorf("Output error %v is not the same as expected error %v", err.Error(), test.errString)
			}
		} else {
			if c != test.c {
				t.Errorf("Output customer %v is not the same as expected customer %v", c, test.c)
			}
		}
	}

}

type convertToCustomersTest struct {
	input     string
	result    map[int]Customer
	errString string
}

var convertToCustomersTests []convertToCustomersTest = []convertToCustomersTest{
	convertToCustomersTest{"", map[int]Customer{}, "Invalid JSON"},
	convertToCustomersTest{"jkhk", map[int]Customer{}, "Invalid JSON"},
	convertToCustomersTest{"{ \"longitude\": 56 }", map[int]Customer{}, "Cannot unmarshal customer"},
	convertToCustomersTest{"{\"latitude\": \"51.92893\", \"user_id\": 1, \"name\": \"Alice Cahill\", \"longitude\": \"-10.27699\"}\n{\"latitude\": \"51.92893\", \"user_id\": 1, \"name\": \"Alice Cahill\", \"longitude\": \"-10.27699\"}",
		map[int]Customer{}, "Customer id overlap"},
	convertToCustomersTest{"{\"latitude\": \"51.92893\", \"user_id\": 1, \"name\": \"Alice Cahill\", \"longitude\": \"-10.27699\"}\n{\"latitude\": \"51.92893\", \"user_id\": 2, \"name\": \"Alice Cahill\", \"longitude\": \"-10.27699\"}",
		map[int]Customer{1: Customer{"51.92893", 1, "Alice Cahill", "-10.27699", greatCircle.Point{greatCircle.DegreeToRadian(-10.27699), greatCircle.DegreeToRadian(51.92893)}},
			2: Customer{"51.92893", 2, "Alice Cahill", "-10.27699", greatCircle.Point{greatCircle.DegreeToRadian(-10.27699), greatCircle.DegreeToRadian(51.92893)}}},
		""},
}

func TestConvertToCustomers(t *testing.T) {
	for _, test := range convertToCustomersTests {
		m, err := convertToCustomers([]byte(test.input))
		if err != nil {
			if !strings.Contains(err.Error(), test.errString) {
				t.Errorf("Output error %v is not the same as expected error %v", err.Error(), test.errString)
			}
		} else {
			if !reflect.DeepEqual(m, test.result) {
				t.Errorf("Output customer %v is not the same as expected customer %v", m, test.result)
			}
		}

	}
}

type getCustomerTest struct {
	method    string
	errString string
	content   string
	fieldName string
	result    string
}

var getCustomerTests []getCustomerTest = []getCustomerTest{
	//Test incorrect method
	getCustomerTest{"POST", "HTTP request is not a PUT request", "", "", ""},
	getCustomerTest{"GET", "HTTP request is not a PUT request", "", "", ""},
	getCustomerTest{"DELETE", "HTTP request is not a PUT request", "", "", ""},
	//Test incorrect field name for the uploaded file
	getCustomerTest{"PUT", "no such file", "", "", ""},
	//Test content validity
	getCustomerTest{"PUT", "Invalid JSON", "cdsc", "customerFile", ""},
	getCustomerTest{"PUT", "Cannot unmarshal customer", "{ \"longitude\": 56 }", "customerFile", ""},
	//Correct content
	//return nothing a customers are too far away (office location is 0,0 by default)
	getCustomerTest{"PUT", "", "{\"latitude\": \"80\", \"user_id\": 1, \"name\": \"user1\", \"longitude\": \"100\"}\n{\"latitude\": \"51.92893\", \"user_id\": 2, \"name\": \"Alice Cahill\", \"longitude\": \"-10.27699\"}", "customerFile", "null"},
	getCustomerTest{"PUT", "", "{\"latitude\": \"0\", \"user_id\": 1, \"name\": \"user1\", \"longitude\": \"0\"}\n{\"latitude\": \"51.92893\", \"user_id\": 2, \"name\": \"Alice Cahill\", \"longitude\": \"-10.27699\"}", "customerFile",
		"[{\"User_id\":1,\"Name\":\"user1\"}]"},
	getCustomerTest{"PUT", "", "{\"latitude\": \"0\", \"user_id\": 1, \"name\": \"user1\", \"longitude\": \"0\"}\n{\"latitude\": \"0\", \"user_id\": 2, \"name\": \"user2\", \"longitude\": \"0\"}", "customerFile",
		"[{\"User_id\":1,\"Name\":\"user1\"},{\"User_id\":2,\"Name\":\"user2\"}]"},
}

func TestGetCustomer(t *testing.T) {
	filePath := "getCustomerTest.txt"

	for _, test := range getCustomerTests {
		body, contentType, err := util.GetByteBuffer(filePath, test.fieldName, test.content)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(test.method, "/v1/customer", body)
		req.Header.Add("Content-Type", contentType)
		writer := httptest.NewRecorder()

		err = GetCustomers(writer, req)

		if err != nil && !strings.Contains(err.Error(), test.errString) {
			t.Errorf("Output error %v is not the same as expected %v", err.Error(), test.errString)
		}
		if result := string(writer.Body.Bytes()); result != test.result {
			t.Errorf("Output result %v is not the same as expected %v", result, test.result)
		}

	}
}
