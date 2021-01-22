package types

import "reflect"

var (
	IdInformationType = reflect.TypeOf(&IdInformation{})
)

type IdInformation struct {
	Id int `db:"id"`
}
