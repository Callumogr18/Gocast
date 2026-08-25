package weathertool

/*
locations.go:
This function simply holds the different provinces that exist in Ireland (4 Total)
and counties (26 total for ROI).
	* ParseProvince(s string) (string, error)
	* ParseCounty(s string)   (string, error)
	* lookup(s string, list []string) (string, error)

lookup is a helper function which simply iterates through the regional and
county lists and checking if the input matches any list items

*/

import (
	"fmt"
	"strings"
)

var Provinces = []string{
	"Connacht",
	"Leinster",
	"Munster",
	"Ulster",
}

var Counties = []string{
	"Antrim",
	"Armagh",
	"Carlow",
	"Cavan",
	"Clare",
	"Cork",
	"Derry",
	"Donegal",
	"Down",
	"Dublin",
	"Fermanagh",
	"Galway",
	"Kerry",
	"Kildare",
	"Kilkenny",
	"Laois",
	"Leitrim",
	"Limerick",
	"Longford",
	"Louth",
	"Mayo",
	"Meath",
	"Monaghan",
	"Offaly",
	"Roscommon",
	"Sligo",
	"Tipperary",
	"Tyrone",
	"Waterford",
	"Westmeath",
	"Wexford",
	"Wicklow",
}

func lookup(s string, list []string) (string, error) {
	s = strings.TrimSpace(s)
	for _, v := range list {
		if strings.EqualFold(s, v) {
			return v, nil
		}
	}

	return "", fmt.Errorf("%w: %q", ErrUnknownLocation, s)
}

func ParseProvince(s string) (string, error) {
	return lookup(s, Provinces)
}

func ParseCounty(s string) (string, error) {
	return lookup(s, Counties)
}
