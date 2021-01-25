package journey

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/nominatim"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var (
	InvalidDetails = errors.New("invalid details")
)

func CreateJourney(httpBody []byte, user *types.User) (*types.Journey, error) {
	journey := &types.Journey{}
	err := json.Unmarshal(httpBody, journey)
	if err != nil {
		return nil, err
	}

	if (journey.Offer && journey.Request) || (!journey.Offer && !journey.Request) || !journey.OpenForRequests {
		return nil, InvalidDetails
	}

	if journey.StartLongLat == nil || journey.EndLongLat == nil {
		return nil, InvalidDetails
	}

	if (journey.TimeIsArrival && journey.TimeIsDeparture) || (!journey.TimeIsArrival && !journey.TimeIsDeparture) {
		return nil, InvalidDetails
	}

	if time.Unix(0, int64(journey.Time*int(time.Millisecond))).Before(time.Now()) {
		return nil, InvalidDetails
	}

	startLongLat := strings.Split(*journey.StartLongLat, ";")
	if len(startLongLat) != 2 {
		return nil, InvalidDetails
	}

	startLong, err := strconv.ParseFloat(startLongLat[0], 32)
	if err != nil {
		return nil, InvalidDetails
	}
	startLat, err := strconv.ParseFloat(startLongLat[1], 32)
	if err != nil {
		return nil, InvalidDetails
	}

	endLongLat := strings.Split(*journey.EndLongLat, ";")
	if len(endLongLat) != 2 {
		return nil, InvalidDetails
	}

	endLong, err := strconv.ParseFloat(endLongLat[0], 32)
	if err != nil {
		return nil, InvalidDetails
	}
	endLat, err := strconv.ParseFloat(endLongLat[1], 32)
	if err != nil {
		return nil, InvalidDetails
	}

	startAddress, err := nominatim.GetAddress(float32(startLong), float32(startLat))
	if err != nil {
		return nil, err
	}
	startAddressString := fmt.Sprintf("%s %s, %s %s", startAddress.Road, startAddress.HouseNumber,
		startAddress.Postcode, startAddress.City)
	journey.StartAddressString = &startAddressString

	endAddress, err := nominatim.GetAddress(float32(endLong), float32(endLat))
	if err != nil {
		return nil, err
	}
	endAddressString := fmt.Sprintf("%s %s, %s %s", endAddress.Road, endAddress.HouseNumber,
		endAddress.Postcode, endAddress.City)
	journey.EndAddressString = &endAddressString

	approximateStartLong := float32(startLong) + 0.001*float32(rand.Intn(3)+1)
	approximateStartLat := float32(startLat) + 0.001*float32(rand.Intn(3)+1)
	approximateEndLong := float32(endLong) + 0.001*float32(rand.Intn(3)+1)
	approximateEndLat := float32(endLat) + 0.001*float32(rand.Intn(3)+1)

	approximateStartLongLat := fmt.Sprintf("%.8g;%.8g", approximateStartLong, approximateStartLat)
	approximateEndLongLat := fmt.Sprintf("%.8g;%.8g", approximateEndLong, approximateEndLat)

	journey.ApproximateStartLongLat = &approximateStartLongLat
	journey.ApproximateEndLongLat = &approximateEndLongLat

	approximateStartAddress, err := nominatim.GetAddress(approximateStartLong, approximateStartLat)
	if err != nil {
		return nil, err
	}
	approximateStartAddressString := fmt.Sprintf("%s %s, %s %s", approximateStartAddress.Road,
		approximateStartAddress.HouseNumber, approximateStartAddress.Postcode, approximateStartAddress.City)
	journey.ApproximateStartAddressString = &approximateStartAddressString

	approximateEndAddress, err := nominatim.GetAddress(approximateEndLong, approximateEndLat)
	if err != nil {
		return nil, err
	}
	approximateEndAddressString := fmt.Sprintf("%s %s, %s %s", approximateEndAddress.Road,
		approximateEndAddress.HouseNumber, approximateEndAddress.Postcode, approximateEndAddress.City)
	journey.ApproximateEndAddressString = &approximateEndAddressString

	journey.Id = nil
	journey.CancelledByAttendeeIds = nil
	journey.CancelledByHostReason = nil
	journey.CancelledByHost = false
	journey.DeclinedUserIds = nil
	journey.AcceptedUserIds = nil
	journey.PendingUserIds = nil
	journey.UserId = user.Id

	return journey, nil
}
