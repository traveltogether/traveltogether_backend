package types

import (
	"github.com/lib/pq"
	"reflect"
)

var (
	JourneyType = reflect.TypeOf(&Journey{})
)

type Journey struct {
	Id                            *int           `json:"id,omitempty" db:"id"`
	UserId                        int            `json:"user_id" db:"user_id"`
	Request                       bool           `json:"request" db:"request"`
	Offer                         bool           `json:"offer" db:"offer"`
	StartLatLong                  *string        `json:"start_lat_long,omitempty" db:"start_lat_long"`
	EndLatLong                    *string        `json:"end_lat_long,omitempty" db:"end_lat_long"`
	ApproximateStartLatLong       *string        `json:"approximate_start_lat_long,omitempty" db:"approximate_start_lat_long"`
	ApproximateEndLatLong         *string        `json:"approximate_end_lat_long,omitempty" db:"approximate_end_lat_long"`
	StartAddressString            *string        `json:"start_address,omitempty" db:"start_address"`
	EndAddressString              *string        `json:"end_address,omitempty" db:"end_address"`
	ApproximateStartAddressString *string        `json:"approximate_start_address,omitempty" db:"approximate_start_address"`
	ApproximateEndAddressString   *string        `json:"approximate_end_address,omitempty" db:"approximate_end_address"`
	Time                          int            `json:"time" db:"time_value"`
	TimeIsDeparture               bool           `json:"time_is_departure" db:"time_is_departure"`
	TimeIsArrival                 bool           `json:"time_is_arrival" db:"time_is_arrival"`
	OpenForRequests               bool           `json:"open_for_requests" db:"open_for_requests"`
	PendingUserIds                *pq.Int64Array `json:"pending_user_ids,omitempty" db:"pending_user_ids"`
	AcceptedUserIds               *pq.Int64Array `json:"accepted_user_ids,omitempty" db:"accepted_user_ids"`
	DeclinedUserIds               *pq.Int64Array `json:"declined_user_ids,omitempty" db:"declined_user_ids"`
	CancelledByHost               bool           `json:"cancelled_by_host" db:"cancelled_by_host"`
	CancelledByHostReason         *string        `json:"cancelled_by_host_reason,omitempty" db:"cancelled_by_host_reason"`
	CancelledByAttendeeIds        *pq.Int64Array `json:"cancelled_by_attendee_ids,omitempty" db:"cancelled_by_attendee_ids"`
	Note                          *string        `json:"note,omitempty" db:"note"`
}
