// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main contains a simple command line tool for Geocoding API
// Directions docs: https://developers.google.com/maps/documentation/geocoding/
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/kr/pretty"
	"google.golang.org/maps"
)

var (
	apiKey       = flag.String("key", "", "API Key for using Google Maps API.")
	address      = flag.String("address", "", "The street address that you want to geocode, in the format used by the national postal service of the country concerned.")
	components   = flag.String("components", "", "A component filter for which you wish to obtain a geocode.")
	bounds       = flag.String("bounds", "", "The bounding box of the viewport within which to bias geocode results more prominently.")
	language     = flag.String("language", "", "The language in which to return results.")
	region       = flag.String("region", "", "The region code, specified as a ccTLD two-character value.")
	latlng       = flag.String("latlng", "", "The textual latitude/longitude value for which you wish to obtain the closest, human-readable address.")
	resultType   = flag.String("result_type", "", "One or more address types, separated by a pipe (|).")
	locationType = flag.String("location_type", "", "One or more location types, separated by a pipe (|).")
)

func usageAndExit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	fmt.Println("Flags:")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Parse()
	client := &http.Client{}
	if *apiKey == "" {
		usageAndExit("Please specify an API Key.")
	}
	ctx := maps.NewContext(*apiKey, client)
	r := &maps.GeocodingRequest{
		Address:  *address,
		Language: *language,
		Region:   *region,
	}

	if *components == "" && *address == "" {
		usageAndExit("Please specify an Address or Components")
	}

	parseComponents(*components, r)
	parseBounds(*bounds, r)
	parseLatLng(*latlng, r)
	parseResultType(*resultType, r)
	parseLocationType(*locationType, r)

	pretty.Println(r)

	resp, err := r.Get(ctx)
	if err != nil {
		log.Fatalf("error %v", err)
	}

	pretty.Println(resp)
}

func parseComponents(components string, r *maps.GeocodingRequest) {
	if components != "" {
		c := strings.Split(components, "|")
		for _, cf := range c {
			i := strings.Split(cf, ":")
			switch {
			case i[0] == "route":
				r.AddComponentFilter(maps.ComponentRoute, i[1])
			case i[0] == "locality":
				r.AddComponentFilter(maps.ComponentLocality, i[1])
			case i[0] == "administrative_area":
				r.AddComponentFilter(maps.ComponentAdministrativeArea, i[1])
			case i[0] == "postal_code":
				r.AddComponentFilter(maps.ComponentPostalCode, i[1])
			case i[0] == "country":
				r.AddComponentFilter(maps.ComponentCounty, i[1])
			}
		}
	}
}

func parseBounds(bounds string, r *maps.GeocodingRequest) {
	if bounds != "" {
		b := strings.Split(bounds, "|")
		sw := strings.Split(b[0], ",")
		ne := strings.Split(b[1], ",")

		swLat, err := strconv.ParseFloat(sw[0], 64)
		if err != nil {
			log.Fatalf("Couldn't parse bounds: %#v", err)
		}
		swLng, err := strconv.ParseFloat(sw[1], 64)
		if err != nil {
			log.Fatalf("Couldn't parse bounds: %#v", err)
		}
		neLat, err := strconv.ParseFloat(ne[0], 64)
		if err != nil {
			log.Fatalf("Couldn't parse bounds: %#v", err)
		}
		neLng, err := strconv.ParseFloat(ne[1], 64)
		if err != nil {
			log.Fatalf("Couldn't parse bounds: %#v", err)
		}

		r.Bounds = &maps.LatLngBounds{
			NorthEast: maps.LatLng{Lat: neLat, Lng: neLng},
			SouthWest: maps.LatLng{Lat: swLat, Lng: swLng},
		}
	}
}

func parseLatLng(latlng string, r *maps.GeocodingRequest) {
	if latlng != "" {
		l := strings.Split(latlng, ",")
		lat, err := strconv.ParseFloat(l[0], 64)
		if err != nil {
			log.Fatalf("Couldn't parse latlng: %#v", err)
		}
		lng, err := strconv.ParseFloat(l[1], 64)
		if err != nil {
			log.Fatalf("Couldn't parse latlng: %#v", err)
		}
		r.LatLng = &maps.LatLng{
			Lat: lat,
			Lng: lng,
		}
	}
}

func parseResultType(resultType string, r *maps.GeocodingRequest) {
	if resultType != "" {
		r.ResultType = strings.Split(resultType, "|")
	}
}

func parseLocationType(locationType string, r *maps.GeocodingRequest) {
	if locationType != "" {
		for _, l := range strings.Split(locationType, "|") {
			switch {
			case l == "ROOFTOP":
				r.LocationType = append(r.LocationType, maps.LocationTypeRooftop)
			case l == "RANGE_INTERPOLATED":
				r.LocationType = append(r.LocationType, maps.LocationTypeRangeInterpolated)
			case l == "GEOMETRIC_CENTER":
				r.LocationType = append(r.LocationType, maps.LocationTypeGeometricCenter)
			case l == "APPROXIMATE":
				r.LocationType = append(r.LocationType, maps.LocationTypeApproximate)
			}
		}

	}
}