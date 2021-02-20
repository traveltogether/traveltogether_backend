package types

import (
	"github.com/lib/pq"
	"reflect"
)

var (
	IdInformationType = reflect.TypeOf(&IdInformation{})
	Int64ArrayType    = reflect.TypeOf(pq.Int64Array{})
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
