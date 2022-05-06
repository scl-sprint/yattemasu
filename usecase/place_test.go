package usecase

import (
	"testing"
	"zemmai-dev/yattemasu/domain/model"
)

func TestMeasureDistance(t *testing.T) {

	want := 8280827.99

	// Tokyo Skytree Coordinates
	loc1 := model.Location{
		Latitude:  35.7100069,
		Longitude: 139.8108103,
	}

	// Sapporo TV Tower Coordinates
	loc2 := model.Location{
		Latitude:  43.061092,
		Longitude: 141.356433,
	}

	have := MeterFromCoordinates(loc1, loc2)

	if int(want) != int(have) {
		t.Fatalf("want %g, but have %g", want, have)
	}
}
