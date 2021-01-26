package helper

import (
	"github.com/forestgiant/sliceutil"
	"github.com/lib/pq"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
)

func ModifyJourneyFields(journey *types.Journey, user *types.User) {
	if journey.UserId != user.Id {
		if journey.PendingUserIds != nil {
			if sliceutil.Contains(*journey.PendingUserIds, int64(user.Id)) {
				journey.PendingUserIds = &pq.Int64Array{int64(user.Id)}

				journey.AcceptedUserIds = nil
				journey.DeclinedUserIds = nil
				journey.StartAddressString = nil
				journey.StartLatLong = nil
				journey.EndAddressString = nil
				journey.EndLatLong = nil
				journey.CancelledByAttendeeIds = nil
				return
			}
		}

		if journey.AcceptedUserIds != nil {
			if sliceutil.Contains(*journey.AcceptedUserIds, int64(user.Id)) {
				journey.AcceptedUserIds = &pq.Int64Array{int64(user.Id)}

				journey.PendingUserIds = nil
				journey.DeclinedUserIds = nil
				journey.CancelledByAttendeeIds = nil
				return
			}
		}

		if journey.DeclinedUserIds != nil {
			if sliceutil.Contains(*journey.DeclinedUserIds, int64(user.Id)) {
				journey.DeclinedUserIds = &pq.Int64Array{int64(user.Id)}

				journey.PendingUserIds = nil
				journey.AcceptedUserIds = nil
				journey.StartAddressString = nil
				journey.StartLatLong = nil
				journey.EndAddressString = nil
				journey.EndLatLong = nil
				journey.CancelledByAttendeeIds = nil
				return
			}
		}

		if journey.CancelledByAttendeeIds != nil {
			if sliceutil.Contains(*journey.CancelledByAttendeeIds, int64(user.Id)) {
				journey.CancelledByAttendeeIds = &pq.Int64Array{int64(user.Id)}

				journey.StartAddressString = nil
				journey.StartLatLong = nil
				journey.EndAddressString = nil
				journey.EndLatLong = nil
				return
			}
		}

		journey.CancelledByAttendeeIds = nil
		journey.PendingUserIds = nil
		journey.AcceptedUserIds = nil
		journey.DeclinedUserIds = nil
		journey.StartAddressString = nil
		journey.StartLatLong = nil
		journey.EndAddressString = nil
		journey.EndLatLong = nil
	}
}
