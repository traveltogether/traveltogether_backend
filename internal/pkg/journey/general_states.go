package journey

import (
	"github.com/lib/pq"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/chat"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/database"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
	"time"
)

func ChangeRequestState(journey *types.Journey, state bool) error {
	if err := database.PrepareAsync(database.DefaultTimeout,
		"UPDATE journeys SET open_for_requests = $1 WHERE id = $2", state, *journey.Id); err != nil {
		return err
	}

	journey.OpenForRequests = state
	return nil
}

func ChangeNote(journey *types.Journey, note *string) error {
	if err := database.PrepareAsync(database.DefaultTimeout,
		"UPDATE journeys SET note = $1 WHERE id = $2", note, *journey.Id); err != nil {
		return err
	}

	journey.Note = note
	return nil
}

func CancelJourney(journey *types.Journey, reason string) error {
	err := database.PrepareAsync(database.DefaultTimeout,
		"UPDATE journeys SET cancelled_by_host = true, cancelled_by_host_reason = $1 WHERE id = $2",
		reason, *journey.Id)

	if err != nil {
		return err
	}

	journey.CancelledByHost = true
	journey.CancelledByHostReason = &reason

	emptySlice := pq.Int64Array{}
	if journey.AcceptedUserIds == nil {
		journey.AcceptedUserIds = &emptySlice
	}
	if journey.PendingUserIds == nil {
		journey.PendingUserIds = &emptySlice
	}

	chatMessage := &types.ChatMessage{}
	chatMessage.Time = int(time.Now().UnixNano() / int64(time.Millisecond))
	chatMessage.Message = reason
	chatMessage.SenderId = &journey.UserId

	for userId := range append(*journey.AcceptedUserIds, *journey.PendingUserIds...) {
		chatMessage.ChatId = nil
		chatMessage.ReceiverId = &userId
		chat.ConnectionManager.SendMessage(chatMessage, journey.UserId, false)
	}

	return nil
}
