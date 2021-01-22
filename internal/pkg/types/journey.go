package types

import "reflect"

var (
	JourneyType = reflect.TypeOf(&Journey{})
)

type Journey struct {
	Id                            *int    `json:"id,omitempty" db:"id"`
	UserId                        int     `json:"user_id" db:"user_id"`
	Request                       bool    `json:"request" db:"request"`
	Offer                         bool    `json:"offer" db:"offer"`
	StartLongLat                  *string `json:"start_long_lat,omitempty" db:"start_long_lat"`
	EndLongLat                    *string `json:"end_long_lat,omitempty" db:"end_long_lat"`
	ApproximateStartLongLat       *string `json:"approximate_start_long_lat,omitempty" db:"approximate_start_long_lat"`
	ApproximateEndLongLat         *string `json:"approximate_end_long_lat,omitempty" db:"approximate_end_long_lat"`
	StartAddressString            *string `json:"start_address,omitempty" db:"start_address"`
	EndAddressString              *string `json:"end_address,omitempty" db:"end_address"`
	ApproximateStartAddressString *string `json:"approximate_start_address,omitempty" db:"approximate_start_address"`
	ApproximateEndAddressString   *string `json:"approximate_end_address,omitempty" db:"approximate_end_address"`
	Time                          int     `json:"time" db:"time_value"`
	TimeIsDeparture               bool    `json:"time_is_departure" db:"time_is_departure"`
	TimeIsArrival                 bool    `json:"time_is_arrival" db:"time_is_arrival"`
	OpenForRequests               bool    `json:"open_for_requests" db:"open_for_requests"`
	PendingUserIds                *[]int  `json:"pending_user_ids,omitempty" db:"pending_user_ids"`
	AcceptedUserIds               *[]int  `json:"accepted_user_ids,omitempty" db:"accepted_user_ids"`
	DeclinedUserIds               *[]int  `json:"declined_user_ids,omitempty" db:"declined_user_ids"`
	CancelledByHost               bool    `json:"cancelled_by_host" db:"cancelled_by_host"`
	CancelledByHostReason         *string `json:"cancelled_by_host_reason" db:"cancelled_by_host_reason"`
	CancelledByAttendeeIds        *[]int  `json:"cancelled_by_attendee_ids" db:"cancelled_by_attendee_ids"`
}
