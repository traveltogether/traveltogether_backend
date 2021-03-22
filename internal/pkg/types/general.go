package types

import (
	"reflect"
)

var (
	IdInformationType = reflect.TypeOf(&IdInformation{})
)

type IdInformation struct {
	Id int `db:"id"`
}

type AddressInformation struct {
	Latitude             float64
	Longitude            float64
	ApproximateLatitude  float64
	ApproximateLongitude float64
	ApproximateLatLong   string
	ApproximateAddress   string
	Try                  int
}
