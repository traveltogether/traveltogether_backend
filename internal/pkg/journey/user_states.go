package journey

import (
	"errors"
	"fmt"
	"github.com/forestgiant/sliceutil"
	"github.com/lib/pq"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/chat"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/database"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/general"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
	"time"
)

var (
	HasBeenCancelled           = errors.New("journey has been cancelled")
	RequestsNotOpen            = errors.New("request not possible")
	UserHasAlreadyBeenDeclined = errors.New("user has already been declined")
	UserHasNotRequestedJoin    = errors.New("user has not requested to join")
	UserHasNotBeenDeclined     = errors.New("user has not been declined")
	UserHasNotBeenAccepted     = errors.New("user has not been accepted")
	AlreadyTookPlace           = errors.New("journey already took place")
	RequestingOwnJourney       = errors.New("requesting own journey")
	UserHasAlreadyBeenAccepted = errors.New("user has already been accepted")
)

func RequestToJoinJourney(journey *types.Journey, userId int64) error {
	if int64(journey.UserId) == userId {
		return RequestingOwnJourney
	}

	if journey.CancelledByHost {
		return HasBeenCancelled
	}

	if !journey.OpenForRequests {
		return RequestsNotOpen
	}

	if time.Unix(0, int64(journey.Time)*int64(time.Millisecond)).Before(time.Now()) {
		return AlreadyTookPlace
	}

	if journey.PendingUserIds != nil {
		if sliceutil.Contains(*journey.PendingUserIds, userId) {
			return nil
		}
	} else {
		journey.PendingUserIds = &pq.Int64Array{}
	}

	if journey.DeclinedUserIds != nil {
		if sliceutil.Contains(*journey.DeclinedUserIds, userId) {
			return UserHasAlreadyBeenDeclined
		}
	}

	if journey.AcceptedUserIds != nil {
		if sliceutil.Contains(*journey.AcceptedUserIds, userId) {
			return UserHasAlreadyBeenAccepted
		}
	}

	if err := database.PrepareAsync(database.DefaultTimeout,
		"UPDATE journeys SET pending_user_ids = array_append(pending_user_ids, $1) "+
			"WHERE id = $2 AND (pending_user_ids IS NULL OR NOT $3 = ANY(pending_user_ids))",
		userId, *journey.Id, userId); err != nil {
		return err
	}

	*journey.PendingUserIds = append(*journey.PendingUserIds, userId)
	return nil
}

func CancelRequestToJoinJourney(journey *types.Journey, userId int64) error {
	if int64(journey.UserId) == userId {
		return RequestingOwnJourney
	}

	if journey.PendingUserIds == nil {
		return UserHasNotBeenAccepted
	}
	if !sliceutil.Contains(*journey.PendingUserIds, userId) {
		return UserHasNotBeenAccepted
	}

	if err := database.PrepareAsync(database.DefaultTimeout, "UPDATE journeys SET pending_user_ids = "+
		"array_remove(pending_user_ids, $1) WHERE id = $2 AND $3=ANY(pending_user_ids)",
		userId, *journey.Id, userId); err != nil {
		return err
	}

	*journey.PendingUserIds = general.RemoveIntFromSlice(*journey.PendingUserIds, userId)
	return nil
}

func AcceptUserToJoinJourney(journey *types.Journey, userId int64) error {
	if journey.CancelledByHost {
		return HasBeenCancelled
	}

	if time.Unix(0, int64(journey.Time)*int64(time.Millisecond)).Before(time.Now()) {
		return AlreadyTookPlace
	}

	if journey.AcceptedUserIds != nil {
		if sliceutil.Contains(*journey.AcceptedUserIds, userId) {
			return nil
		}
	} else {
		journey.AcceptedUserIds = &pq.Int64Array{}
	}

	if journey.DeclinedUserIds != nil {
		if sliceutil.Contains(*journey.DeclinedUserIds, userId) {
			return UserHasAlreadyBeenDeclined
		}
	}

	if journey.PendingUserIds == nil {
		return UserHasNotRequestedJoin
	}
	if !sliceutil.Contains(*journey.PendingUserIds, userId) {
		return UserHasNotRequestedJoin
	}

	if err := database.PrepareAsync(database.DefaultTimeout, "UPDATE journeys SET accepted_user_ids = "+
		"array_append(accepted_user_ids, $1) WHERE id = $2 "+
		"AND (accepted_user_ids IS NULL OR NOT $3 = ANY(accepted_user_ids))",
		userId, *journey.Id, userId); err != nil {
		return err
	}

	*journey.AcceptedUserIds = append(*journey.AcceptedUserIds, userId)

	if err := database.PrepareAsync(database.DefaultTimeout, "UPDATE journeys SET pending_user_ids = "+
		"array_remove(pending_user_ids, $1) WHERE id = $2 AND $3 = ANY(pending_user_ids)",
		userId, *journey.Id, userId); err != nil {
		return err
	}

	*journey.PendingUserIds = general.RemoveIntFromSlice(*journey.PendingUserIds, userId)
	return nil
}

func CancelAcceptToJoinJourney(journey *types.Journey, userId int64) error {
	if journey.AcceptedUserIds == nil {
		return UserHasNotBeenAccepted
	}
	if !sliceutil.Contains(*journey.AcceptedUserIds, userId) {
		return UserHasNotBeenAccepted
	}

	if err := database.PrepareAsync(database.DefaultTimeout, "UPDATE journeys SET accepted_user_ids = "+
		"array_remove(accepted_user_ids, $1) WHERE id = $2 AND $3 = ANY(accepted_user_ids)",
		userId, *journey.Id, userId); err != nil {
		return err
	}

	*journey.PendingUserIds = general.RemoveIntFromSlice(*journey.PendingUserIds, userId)
	return nil
}

func DeclineUserToJoinJourney(journey *types.Journey, userId int64) error {
	if journey.CancelledByHost {
		return HasBeenCancelled
	}

	if time.Unix(0, int64(journey.Time)*int64(time.Millisecond)).Before(time.Now()) {
		return AlreadyTookPlace
	}

	if journey.PendingUserIds == nil {
		return UserHasNotRequestedJoin
	}
	if !sliceutil.Contains(*journey.PendingUserIds, userId) {
		return UserHasNotRequestedJoin
	}

	if journey.AcceptedUserIds != nil {
		if sliceutil.Contains(*journey.AcceptedUserIds, userId) {
			return UserHasAlreadyBeenAccepted
		}
	}

	if journey.DeclinedUserIds != nil {
		if sliceutil.Contains(*journey.DeclinedUserIds, userId) {
			return nil
		}
	} else {
		journey.DeclinedUserIds = &pq.Int64Array{}
	}

	if err := database.PrepareAsync(database.DefaultTimeout, "UPDATE journeys SET declined_user_ids = "+
		"array_append(declined_user_ids, $1) WHERE id = $2 "+
		"AND (declined_user_ids IS NULL OR NOT $3 = ANY(declined_user_ids))",
		userId, *journey.Id, userId); err != nil {
		return err
	}

	*journey.DeclinedUserIds = append(*journey.DeclinedUserIds, userId)

	if err := database.PrepareAsync(database.DefaultTimeout, "UPDATE journeys SET pending_user_ids = "+
		"array_remove(pending_user_ids, $1) WHERE id = $2 AND $3 = ANY(pending_user_ids)",
		userId, *journey.Id, userId); err != nil {
		return err
	}

	*journey.PendingUserIds = general.RemoveIntFromSlice(*journey.PendingUserIds, userId)
	return nil
}

func ReverseDeclineUserToJoinJourney(journey *types.Journey, userId int64) error {
	if journey.CancelledByHost {
		return HasBeenCancelled
	}

	if time.Unix(0, int64(journey.Time)*int64(time.Millisecond)).Before(time.Now()) {
		return AlreadyTookPlace
	}

	if journey.DeclinedUserIds == nil {
		return UserHasNotBeenDeclined
	}
	if !sliceutil.Contains(*journey.DeclinedUserIds, userId) {
		return UserHasNotBeenDeclined
	}

	if err := database.PrepareAsync(database.DefaultTimeout, "UPDATE journeys SET declined_user_ids = "+
		"array_remove(declined_user_ids, $1) WHERE id = $2 AND $3 = ANY(declined_user_ids)",
		userId, *journey.Id, userId); err != nil {
		return err
	}

	*journey.DeclinedUserIds = general.RemoveIntFromSlice(*journey.DeclinedUserIds, userId)
	return nil
}

func CancelAttendanceAtJourney(journey *types.Journey, user *types.User) error {
	if int64(journey.UserId) == int64(user.Id) {
		return RequestingOwnJourney
	}

	if journey.AcceptedUserIds == nil {
		return UserHasNotBeenAccepted
	}
	if !sliceutil.Contains(*journey.AcceptedUserIds, int64(user.Id)) {
		return UserHasNotBeenAccepted
	}

	if journey.CancelledByAttendeeIds == nil {
		journey.CancelledByAttendeeIds = &pq.Int64Array{}
	}

	if err := database.PrepareAsync(database.DefaultTimeout, "UPDATE journeys SET cancelled_by_attendee_ids = "+
		"array_append(cancelled_by_attendee_ids, $1) WHERE id = $2 "+
		"AND (cancelled_by_attendee_ids IS NULL OR NOT $3=ANY(cancelled_by_attendee_ids))",
		int64(user.Id), *journey.Id, int64(user.Id)); err != nil {
		return err
	}

	chatMessage := &types.ChatMessage{}
	chatMessage.Time = int(time.Now().UnixNano() / int64(time.Millisecond))
	chatMessage.Message = fmt.Sprintf("User %s left the journey.", user.Username)
	chatMessage.SenderId = &user.Id
	chatMessage.ReceiverId = &journey.UserId

	chat.ConnectionManager.SendMessage(chatMessage, user.Id, false)

	*journey.CancelledByAttendeeIds = append(*journey.CancelledByAttendeeIds, int64(user.Id))

	if err := database.PrepareAsync(database.DefaultTimeout, "UPDATE journeys SET accepted_user_ids = "+
		"array_remove(accepted_user_ids, $1) WHERE id = $2 AND $3 = ANY(accepted_user_ids)",
		int64(user.Id), *journey.Id, int64(user.Id)); err != nil {
		return err
	}

	*journey.AcceptedUserIds = general.RemoveIntFromSlice(*journey.AcceptedUserIds, int64(user.Id))
	return nil
}
