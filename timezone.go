package p24

import (
	_ "embed"
	"time"
)

// timeZoneEuropeKiev holds Europe/Kiev in the IANA Time Zone database-format
//go:embed timezone/Kiev
var timeZoneEuropeKiev []byte

var kievLocation = NewKievLocation()

// NewKievLocation returns time.Location of Europe/Kiev time zone
func NewKievLocation() *time.Location {
	l, _ := time.LoadLocationFromTZData("Europe/Kiev", timeZoneEuropeKiev)
	return l
}
