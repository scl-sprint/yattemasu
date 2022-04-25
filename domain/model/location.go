package model

import "fmt"

type Location struct {
	Longitude float64
	Latitude float64
}

func (loc Location) String() string {
	return fmt.Sprintf("%f,%f", loc.Longitude, loc.Latitude)
}