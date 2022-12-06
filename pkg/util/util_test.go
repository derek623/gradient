package util

import (
	"net/http/httptest"
	"strings"
	"testing"
)

type equalTest struct {
	left, right float64
	expected    bool
}

var equalTests []equalTest = []equalTest{
	equalTest{0.0, 0.0, true},
	equalTest{0.0, 1e-9, true},
	equalTest{0.0, -1e-9, true},
	equalTest{1e-9, 0.0, true},
	equalTest{-1e-9, 0.0, true},
	equalTest{1.0, 1 + -1e-9, true},
	equalTest{-1.0, -1 - 1e-9, true},
	equalTest{400.01, 400.01 + -1e-9, true},
	equalTest{-400.01, -400.01 - 1e-9, true},

	equalTest{0.0, 1e-8, false},
	equalTest{1e-8, 0.0, false},
	equalTest{1 + 1e-7, 1.0, false},
	equalTest{-1 - 1e-7, -1.0, false},
	equalTest{1.0, 1 + 1e-7, false},
	equalTest{1.0, -1 - 1e-7, false},
	equalTest{400.01, 400.01 + 1e-7, false},
	equalTest{400.01, 400.01 - 1e-7, false},
}

func TestEqual(t *testing.T) {
	for _, test := range equalTests {
		if r := Equal(test.left, test.right); r != test.expected {
			t.Errorf("Output %v not equal to expected %v", r, test.expected)
		}
	}
}

type smallerOrEqualTest struct {
	left, right float64
	expected    bool
}

var smallerOrEqualTests []smallerOrEqualTest = []smallerOrEqualTest{
	smallerOrEqualTest{0.0, 0.0, true},
	smallerOrEqualTest{0.0, 1e-9, true},
	smallerOrEqualTest{0.0, -1e-9, true},
	smallerOrEqualTest{1e-9, 0.0, true},
	smallerOrEqualTest{-1e-9, 0.0, true},
	smallerOrEqualTest{1.0, 1 + -1e-9, true},
	smallerOrEqualTest{-1.0, -1 - 1e-9, true},
	smallerOrEqualTest{400.01, 400.01 + -1e-9, true},
	smallerOrEqualTest{-400.01, -400.01 - 1e-9, true},

	smallerOrEqualTest{0.0, 1e-8, true},
	smallerOrEqualTest{1e-8, 0.0, false},
	smallerOrEqualTest{1 + 1e-7, 1.0, false},
	smallerOrEqualTest{1.0, 1 + 1e-7, true},
	smallerOrEqualTest{400.01 + 1e-7, 400.01, false},
	smallerOrEqualTest{400.01, 400.01 + 1e-7, true},
}

func TestSmallerOrEqual(t *testing.T) {
	for _, test := range smallerOrEqualTests {
		if r := SmallerOrEqual(test.left, test.right); r != test.expected {
			t.Errorf("Output %v not equal to expected %v", r, test.expected)
		}
	}
}

type largerOrEqualTest struct {
	left, right float64
	expected    bool
}

var largerOrEqualTests []largerOrEqualTest = []largerOrEqualTest{
	largerOrEqualTest{0.0, 0.0, true},
	largerOrEqualTest{0.0, 1e-9, true},
	largerOrEqualTest{0.0, -1e-9, true},
	largerOrEqualTest{1e-9, 0.0, true},
	largerOrEqualTest{-1e-9, 0.0, true},
	largerOrEqualTest{1.0, 1 + -1e-9, true},
	largerOrEqualTest{-1.0, -1 - 1e-9, true},
	largerOrEqualTest{400.01, 400.01 + -1e-9, true},
	largerOrEqualTest{-400.01, -400.01 - 1e-9, true},

	largerOrEqualTest{0.0, 1e-8, false},
	largerOrEqualTest{1e-8, 0.0, true},
	largerOrEqualTest{1 + 1e-7, 1.0, true},
	largerOrEqualTest{1.0, 1 + 1e-7, false},
	largerOrEqualTest{400.01 + 1e-7, 400.01, true},
	largerOrEqualTest{400.01, 400.01 + 1e-7, false},
}

func TestLargerOrEqual(t *testing.T) {
	for _, test := range largerOrEqualTests {
		if r := LargerOrEqual(test.left, test.right); r != test.expected {
			t.Errorf("Output %v not equal to expected %v", r, test.expected)
		}
	}
}

type getFileInBytesTest struct {
	errString string
	content   string
	fieldName string
	result    string
}

var getFileInBytesTests []getFileInBytesTest = []getFileInBytesTest{
	getFileInBytesTest{"no such file", "", "", ""},
	getFileInBytesTest{"no such file", "", "fff", ""},
	getFileInBytesTest{"", "", "customerFile", ""},
	getFileInBytesTest{"", "", "customerFile", ""},
	getFileInBytesTest{"", "a", "customerFile", "a"},
	getFileInBytesTest{"", "{}", "customerFile", "{}"},
	getFileInBytesTest{"", "{\"value1\": 1, \"value2\": \"rtr\"}", "customerFile", "{\"value1\": 1, \"value2\": \"rtr\"}"},
}

func TestGetFileInBytes(t *testing.T) {
	filePath := "getFileInBytesTest.txt"

	for _, test := range getFileInBytesTests {
		body, contentType, err := GetByteBuffer(filePath, test.fieldName, test.content)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest("", "/test", body)
		req.Header.Add("Content-Type", contentType)

		b, er := GetFileInBytes(req)

		if er != nil && !strings.Contains(er.Error(), test.errString) {
			t.Errorf("Output error %v is not the same as expected %v", er.Error(), test.errString)
		}
		if result := string(b); result != test.result {
			t.Errorf("Output result %v is not the same as expected %v", result, test.result)
		}

	}
}
