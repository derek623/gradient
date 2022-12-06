// Package greatCircle provide functionalities for calculating the great circle distance and generating Point on a sphere
package greatCircle

import (
	"fmt"
	"math"

	"git.codesubmit.io/sfox/party-invite-ruiegv/pkg/util"
)

const Radius = 6371.009

// Input for the Sin, Cos, ASin... function are float64
type Point struct {
	Longitude float64
	Latitude  float64
}

// Implement stringer for printing the Point struct properly
func (p Point) String() string {
	return fmt.Sprintf("{Longitude: %f, Latitude: %f}", p.Longitude, p.Latitude)
}

// Check if the point is valid
func (p Point) Valid() bool {
	return validLongtitude(p.Longitude) && validLatitude(p.Latitude)
}

// Check if the provided longitude is valid
func validLongtitude(longitude float64) bool {
	return util.LargerOrEqual(longitude, -math.Pi) && util.SmallerOrEqual(longitude, math.Pi)
}

// Check if the provided latitude is valid
func validLatitude(latitude float64) bool {
	return util.LargerOrEqual(latitude, -math.Pi/2.0) && util.SmallerOrEqual(latitude, math.Pi/2.0)
}

// Generate a Point struct with the provided longitude and latitude
func MakePoint(longitude float64, latitude float64) Point {
	return Point{longitude, latitude}
}

// Helper to convert degree to radian
func DegreeToRadian(degree float64) float64 {
	return degree * math.Pi / 180.0
}

// Return the distance between 2 points. Assume longtitude and latitdue are in radian already
func Distance(p1 Point, p2 Point, radius float64) float64 {
	centralAngle := math.Acos(math.Sin(p1.Latitude)*math.Sin(p2.Latitude) + math.Cos(p1.Latitude)*math.Cos(p2.Latitude)*math.Cos(math.Abs(p1.Longitude-p2.Longitude)))
	return radius * centralAngle
}
