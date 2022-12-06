package greatCircle

import (
	"math"
	"testing"

	"git.codesubmit.io/sfox/party-invite-ruiegv/pkg/util"
)

type makePointTest struct {
	longitude, latitude float64
	expected            Point
}

var makePointTests = []makePointTest{
	//both 0
	makePointTest{0, 0, Point{0, 0}},
	//only longitude is 0
	makePointTest{0, 1.0, Point{0, 1.0}},
	makePointTest{0, 145.5, Point{0, 145.5}},
	makePointTest{0, -1.0, Point{0, -1.0}},
	makePointTest{0, -346.0, Point{0, -346.0}},
	//only latitude is 0
	makePointTest{1.0, 0, Point{1.0, 0}},
	makePointTest{0.56, 0, Point{0.56, 0}},
	makePointTest{-1.0, 0, Point{-1.0, 0}},
	makePointTest{-1.67, 0, Point{-1.67, 0}},
	//Both positive
	makePointTest{0.1, 0.1, Point{0.1, 0.1}},
	makePointTest{100.87, 5.1, Point{100.87, 5.1}},
	//Both negative
	makePointTest{-0.1, -0.1, Point{-0.1, -0.1}},
	makePointTest{-100.1, -3.1, Point{-100.1, -3.1}},
}

func TestMakePoint(t *testing.T) {
	for _, test := range makePointTests {
		if p := MakePoint(test.longitude, test.latitude); p != test.expected {
			t.Errorf("Output %v not equal to expected %v", p, test.expected)
		}

	}
}

type degreeToRadianTest struct {
	degree, expected float64
}

var degreeToRadianTests []degreeToRadianTest = []degreeToRadianTest{
	degreeToRadianTest{0.0, 0.0},
	//well none values
	degreeToRadianTest{180, math.Pi},
	degreeToRadianTest{360, math.Pi * 2},
	degreeToRadianTest{-180, -math.Pi},
	degreeToRadianTest{-360, -math.Pi * 2},
}

func TestDegreeToRadian(t *testing.T) {
	for _, test := range degreeToRadianTests {
		if r := DegreeToRadian(test.degree); r != test.expected {
			t.Errorf("Output %v not equal to expected %v", r, test.expected)
		}
	}
}

type distanceTest struct {
	p1, p2   Point
	radius   float64
	expected float64
}

var distanceTests []distanceTest = []distanceTest{
	distanceTest{Point{DegreeToRadian(0), DegreeToRadian(0)}, Point{DegreeToRadian(0), DegreeToRadian(0)}, 10, 0},
	distanceTest{Point{DegreeToRadian(28), DegreeToRadian(55)}, Point{DegreeToRadian(28), DegreeToRadian(55)}, 10, 0},
	distanceTest{Point{DegreeToRadian(28), DegreeToRadian(55)}, Point{DegreeToRadian(28), DegreeToRadian(55)}, 6371.009, 0},
	distanceTest{Point{DegreeToRadian(-6.257664), DegreeToRadian(53.339428)}, Point{DegreeToRadian(-6.257664), DegreeToRadian(53.339428)}, 6371.009, 0},
	distanceTest{Point{DegreeToRadian(-7.257664), DegreeToRadian(53.339428)}, Point{DegreeToRadian(-6.257664), DegreeToRadian(53.339428)}, 6371.009, 66.391069412779},
	distanceTest{Point{DegreeToRadian(28), DegreeToRadian(55)}, Point{DegreeToRadian(70), DegreeToRadian(86)}, 10, 5.60686309},
}

func TestDistance(t *testing.T) {
	for _, test := range distanceTests {
		if d := Distance(test.p1, test.p2, test.radius); !util.Equal(d, test.expected) {
			t.Errorf("Output %v not equal to expected %v", d, test.expected)
		}
	}
}

type validTest struct {
	longtitude, latitude float64
	expected             bool
}

var validTests []validTest = []validTest{
	validTest{0, 0, true},
	//check latitude
	validTest{0, math.Pi / 4, true},
	validTest{0, math.Pi / 2, true},
	validTest{0, -math.Pi / 4, true},
	validTest{0, -math.Pi / 4, true},
	validTest{0, math.Pi/2 + 0.2, false},
	validTest{0, -math.Pi/2 - 0.2, false},
	//Check longitude
	validTest{math.Pi, 0, true},
	validTest{math.Pi / 2, 0, true},
	validTest{-math.Pi, 0, true},
	validTest{-math.Pi / 2, 0, true},
	validTest{math.Pi + 0.1, 0, false},
	validTest{-math.Pi - 0.1, 0, false},
	//Check both
	validTest{math.Pi, math.Pi / 2, true},
	validTest{-math.Pi, -math.Pi / 2, true},
	validTest{-math.Pi, math.Pi / 2, true},
	validTest{math.Pi, -math.Pi / 2, true},
	validTest{math.Pi, -math.Pi/2 - 0.1, false},
	validTest{math.Pi + 0.2, -math.Pi / 2, false},
	validTest{math.Pi + 0.2, -math.Pi/2 - 0.1, false},
}

func TestValid(t *testing.T) {
	for _, test := range validTests {
		p := MakePoint(test.longtitude, test.latitude)
		if p.Valid() != test.expected {
			t.Errorf("Output %v not equal to expected %v", p.Valid(), test.expected)
		}
	}
}
