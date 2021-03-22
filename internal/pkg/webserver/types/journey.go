package types

import "github.com/traveltogether/traveltogether_backend/internal/pkg/types"

type Journeys struct {
	Journeys []*types.Journey `json:"journeys"`
}
