package model

import (
	"fmt"
	"os"

	"googlemaps.github.io/maps"
)

type Place struct {
	Name     string
	Location Location
	ImageUrl string
	Address  string
}

func (p *Place) FromSearchResult(result maps.PlacesSearchResult) *Place {
	p.Name = result.Name
	if len(result.Photos) != 0 {
		p.ImageUrl = fmt.Sprintf("https://maps.googleapis.com/maps/api/place/photo?maxwidth=400&photoreference=%s&key=%s", result.Photos[0].PhotoReference, os.Getenv("GMAP_API_KEY"))
	} else {
		p.ImageUrl = "https://placehold.jp/400x300.png"
	}
	p.Address = result.FormattedAddress
	return p
}
