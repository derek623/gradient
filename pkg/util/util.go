package util

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"os"
)

// This is for floating point comparision.
const Epsilon = 1e-8

// Test if 2 floating point numbers are equal
func Equal(num1 float64, num2 float64) bool {
	return math.Abs(num1-num2) < Epsilon
}

// Test if floating point num1 is < floating point num2
func SmallerOrEqual(num1 float64, num2 float64) bool {
	return num1 < num2 || Equal(num1, num2)
}

// Test if floating point num1 is > floating point num2
func LargerOrEqual(num1 float64, num2 float64) bool {
	return num1 > num2 || Equal(num1, num2)
}

// Read the uploaded file from a http request and return it in a byte array
func GetFileInBytes(request *http.Request) ([]byte, error) {
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	request.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, _, err := request.FormFile("customerFile")
	if err != nil {
		//log.Println(err)
		return nil, err
	}
	defer file.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}

// A generic handler for http request
func ErrorHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, err.Error())
			log.Println(err.Error())
		}
	}
}

// A helper function for unit testing to generate byte buffer
func GetByteBuffer(filePath string, fieldName string, content string) (*bytes.Buffer, string, error) {
	c := []byte(content)
	//write test content to file filePath
	e := os.WriteFile(filePath, c, 0644)
	if e != nil {
		return nil, "", e
	}

	//Create a form file for the request
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", err
	}

	w, err := mw.CreateFormFile(fieldName, filePath)
	if err != nil {
		return nil, "", err
	}

	if _, err := io.Copy(w, file); err != nil {
		return nil, "", err
	}

	// close the writer before making the request
	file.Close()
	mw.Close()

	return body, mw.FormDataContentType(), nil
}
