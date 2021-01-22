package journey

import (
	"errors"
	"github.com/forestgiant/sliceutil"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/database"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/general"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
)

var (
	HasBeenCancelled           = errors.New("journey has been cancelled")
	RequestsNotOpen            = errors.New("request not possible")
	UserHasAlreadyBeenDeclined = errors.New("user has already been declined")
	UserHasNotRequestedJoin    = errors.New("user has not requested to join")
	UserHasNotBeenDeclined     = errors.New("user has not been declined")
	UserHasNotBeenAccepted     = errors.New("user has not been accepted")
)

func RequestToJoinJourney(journey *types.Journey, userId int) error {
	if journey.CancelledByHost {
		return HasBeenCancelled
	}

	if !journey.OpenForRequests {
		return RequestsNotOpen
	}

	if journey.DeclinedUserIds != nil {
		if sliceutil.Contains(*journey.DeclinedUserIds, userId) {
			return UserHasAlreadyBeenDeclined
		}
	}

	if journey.PendingUserIds != nil {
		if sliceutil.Contains(*journey.PendingUserIds, userId) {
			return nil
		}
	} else {
		*journey.PendingUserIds = []int{}
	}

	if err := database.PrepareAsync(database.DefaultTimeout,
		"UPDATE journeys SET pending_user_ids = array_append(pending_user_ids, $1) "+
			"WHERE id = $2 AND NOT (pending_user_ids @> ARRAY[$3])", userId, *journey.Id, userId); err != nil {
		return err
	}

	*journey.PendingUserIds = append(*journey.PendingUserIds, userId)
	return nil
}

func AcceptUserToJoinJourney(journey *types.Journey, userId int) error {
	if journey.CancelledByHost {
		return HasBeenCancelled
	}

	if journey.PendingUserIds == nil {
		return UserHasNotRequestedJoin
	}
	if !sliceutil.Compare(*journey.PendingUserIds, userId) {
		return UserHasNotRequestedJoin
	}

	if journey.AcceptedUserIds != nil {
		if sliceutil.Contains(*journey.AcceptedUserIds, userId) {
			return nil
		}
	} else {
		*journey.AcceptedUserIds = []int{}
	}

	if err := database.PrepareAsync(database.DefaultTimeout, "UPDATE journeys SET accepted_user_ids = "+
		"array_append(accepted_user_ids, $1) WHERE id = $2 AND NOT (accepted_user_ids @> ARRAY[$3])",
		userId, *journey.Id, userId); err != nil {
		return err
	}

	*journey.AcceptedUserIds = append(*journey.AcceptedUserIds, userId)

	if err := database.PrepareAsync(database.DefaultTimeout, "UPDATE journeys SET pending_user_ids = "+
		"array_remove(pending_user_ids, $1) WHERE id = $2 AND (pending_user_ids @> ARRAY[$3])",
		userId, *journey.Id, userId); err != nil {
		return err
	}

	*journey.PendingUserIds = general.RemoveIntFromSlice(*journey.PendingUserIds, userId)
	return nil
}

func DeclineUserToJoinJourney(journey *types.Journey, userId int) error {
	if journey.PendingUserIds == nil {
		return UserHasNotRequestedJoin
	}
	if !sliceutil.Contains(*journey.PendingUserIds, userId) {
		return UserHasNotRequestedJoin
	}

	if journey.DeclinedUserIds != nil {
		if sliceutil.Contains(*journey.DeclinedUserIds, userId) {
			return nil
		}
	} else {
		*journey.DeclinedUserIds = []int{}
	}

	if err := database.PrepareAsync(database.DefaultTimeout, "UPDATE journeys SET declined_user_ids = "+
		"array_append(declined_user_ids, $1) WHERE id = $2 AND NOT (declined_user_ids @> ARRAY[$3])",
		userId, *journey.Id, userId); err != nil {
		return err
	}

	*journey.DeclinedUserIds = append(*journey.DeclinedUserIds, userId)

	if err := database.PrepareAsync(database.DefaultTimeout, "UPDATE journeys SET pending_user_ids = "+
		"array_remove(pending_user_ids, $1) WHERE id = $2 AND (pending_user_ids @> ARRAY[$3])",
		userId, *journey.Id, userId); err != nil {
		return err
	}

	*journey.PendingUserIds = general.RemoveIntFromSlice(*journey.PendingUserIds, userId)
	return nil
}

func ReverseDeclineUserToJoinJourney(journey *types.Journey, userId int) error {
	if journey.DeclinedUserIds == nil {
		return UserHasNotBeenDeclined
	}
	if !sliceutil.Contains(*journey.DeclinedUserIds, userId) {
		return UserHasNotBeenDeclined
	}

	if err := database.PrepareAsync(database.DefaultTimeout, "UPDATE journeys SET declined_user_ids = "+
		"array_remove(declined_user_ids, $1) WHERE id = $2 AND (declined_user_ids @> ARRAY[$3])",
		userId, *journey.Id, userId); err != nil {
		return err
	}

	*journey.DeclinedUserIds = general.RemoveIntFromSlice(*journey.DeclinedUserIds, userId)
	return nil
}

func CancelAttendanceAtJourney(journey *types.Journey, userId int) error {
	if journey.AcceptedUserIds == nil {
		return UserHasNotBeenAccepted
	}
	if !sliceutil.Contains(*journey.AcceptedUserIds, userId) {
		return UserHasNotBeenAccepted
	}

	if journey.CancelledByAttendeeIds == nil {
		*journey.CancelledByAttendeeIds = []int{}
	}

	if err := database.PrepareAsync(database.DefaultTimeout, "UPDATE journeys SET cancelled_by_attendee_ids = "+
		"array_append(cancelled_by_attendee_ids, $1) WHERE id = $2 AND (cancelled_by_attendee_ids @> ARRAY[$3])",
		userId, *journey.Id, userId); err != nil {
		return err
	}

	// TODO notify host about change

	*journey.CancelledByAttendeeIds = append(*journey.CancelledByAttendeeIds, userId)

	if err := database.PrepareAsync(database.DefaultTimeout, "UPDATE journeys SET accepted_user_ids = "+
		"array_remove(accepted_user_ids, $1) WHERE id = $2 AND (accepted_user_ids @> ARRAY[$3])",
		userId, *journey.Id, userId); err != nil {
		return err
	}

	*journey.AcceptedUserIds = general.RemoveIntFromSlice(*journey.AcceptedUserIds, userId)
	return nil
}
