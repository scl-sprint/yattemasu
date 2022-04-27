package usecase

import (
	"math"
	"zemmai-dev/yattemasu/domain/model"
)

func MeterFromCoordinates(loc1 model.Location, loc2 model.Location) float64 {

	// code reference: https://qiita.com/YasumiYasumi/items/3ed8ea69b85dac055381

	const A = 6371008
	const B = 6371008

	const F = (A - B) / A

	lat1 := loc1.Latitude * math.Pi / 180
	lat2 := loc2.Latitude * math.Pi / 180
	lng1 := loc1.Longitude * math.Pi / 180
	lng2 := loc2.Longitude * math.Pi / 180

	phi1 := math.Atan(B / A * math.Tan(lat1))
	phi2 := math.Atan(B / A * math.Tan(lat2))

	f1 := math.Sin(phi1) * math.Sin(phi2)
	f2 := math.Cos(phi1) * math.Cos(phi2)
	f3 := math.Cos(lng1 - lng2)

	X := math.Acos(f1 + f2*f3)

	f4 := math.Sin(X) - X
	f5 := math.Sin(phi1) + math.Sin(phi2)
	f6 := math.Cos(X/2) * math.Cos(X/2)

	f7 := math.Sin(X) + X
	f8 := math.Sin(phi1) - math.Sin(phi2)
	f9 := math.Sin(X/2) * math.Sin(X/2)

	drho := F / 8 * (f4*f5*f5/f6 - f7*f8*f8/f9)

	meter := A * (X + drho)

	return meter
}
