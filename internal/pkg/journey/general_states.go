package journey

import (
	"github.com/traveltogether/traveltogether_backend/internal/pkg/database"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
)

func ChangeRequestState(journey *types.Journey, state bool) error {
	if err := database.PrepareAsync(database.DefaultTimeout,
		"UPDATE journeys SET open_for_requests = $1 WHERE id = $2", state, *journey.Id); err != nil {
		return err
	}

	journey.OpenForRequests = state
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
	*journey.CancelledByHostReason = reason

	// TODO notify user via chat about cancelled journey

	return nil
}
