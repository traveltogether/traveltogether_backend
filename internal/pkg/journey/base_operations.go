package journey

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/nominatim"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
	"math/big"
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

	if (journey.Offer && journey.Request) || (!journey.Offer && !journey.Request) {
		return nil, InvalidDetails
	}

	if journey.StartLatLong == nil || journey.EndLatLong == nil {
		return nil, InvalidDetails
	}

	if journey.StartLatLong == journey.EndLatLong {
		return nil, InvalidDetails
	}

	if (journey.TimeIsArrival && journey.TimeIsDeparture) || (!journey.TimeIsArrival && !journey.TimeIsDeparture) {
		return nil, InvalidDetails
	}

	if time.Unix(0, int64(journey.Time*int(time.Millisecond))).Before(time.Now()) {
		return nil, InvalidDetails
	}

	startLatLong := strings.Split(*journey.StartLatLong, ";")
	if len(startLatLong) != 2 {
		return nil, InvalidDetails
	}

	startLat, err := strconv.ParseFloat(startLatLong[0], 32)
	if err != nil {
		return nil, InvalidDetails
	}
	startLong, err := strconv.ParseFloat(startLatLong[1], 32)
	if err != nil {
		return nil, InvalidDetails
	}

	endLatLong := strings.Split(*journey.EndLatLong, ";")
	if len(endLatLong) != 2 {
		return nil, InvalidDetails
	}

	endLat, err := strconv.ParseFloat(endLatLong[0], 32)
	if err != nil {
		return nil, InvalidDetails
	}
	endLong, err := strconv.ParseFloat(endLatLong[1], 32)
	if err != nil {
		return nil, InvalidDetails
	}

	startAddress, err := nominatim.GetAddress(float32(startLat), float32(startLong))
	if err != nil {
		return nil, err
	}
	startAddressString := fmt.Sprintf("%s %s, %s %s", startAddress.Road, startAddress.HouseNumber,
		startAddress.Postcode, startAddress.City)
	journey.StartAddressString = &startAddressString

	endAddress, err := nominatim.GetAddress(float32(endLat), float32(endLong))
	if err != nil {
		return nil, err
	}
	endAddressString := fmt.Sprintf("%s %s, %s %s", endAddress.Road, endAddress.HouseNumber,
		endAddress.Postcode, endAddress.City)
	journey.EndAddressString = &endAddressString

	if startAddressString == endAddressString {
		return nil, InvalidDetails
	}

	approximateStartAddressInformation := &types.AddressInformation{
		Latitude:  startLat,
		Longitude: startLong,
	}

	err = getApproximateAddress(approximateStartAddressInformation)
	if err != nil {
		return nil, err
	}

	journey.ApproximateStartLatLong = &approximateStartAddressInformation.ApproximateLatLong
	journey.ApproximateStartAddressString = &approximateStartAddressInformation.ApproximateAddress

	approximateEndAddressInformation := &types.AddressInformation{
		Latitude:  endLat,
		Longitude: endLong,
	}

	err = getApproximateAddress(approximateEndAddressInformation)
	if err != nil {
		return nil, err
	}

	journey.ApproximateEndLatLong = &approximateEndAddressInformation.ApproximateLatLong
	journey.ApproximateEndAddressString = &approximateEndAddressInformation.ApproximateAddress

	journey.Id = nil
	journey.CancelledByAttendeeIds = nil
	journey.CancelledByHostReason = nil
	journey.CancelledByHost = false
	journey.DeclinedUserIds = nil
	journey.AcceptedUserIds = nil
	journey.PendingUserIds = nil
	journey.UserId = user.Id
	journey.OpenForRequests = true

	return journey, nil
}

func getApproximateAddress(addressInformation *types.AddressInformation) error {
	addressInformation.Try += 1

	if addressInformation.Try == 3 {
		return nil
	}

	approximateLat, err := offsetCoordinate(addressInformation.Latitude)
	if err != nil {
		return err
	}
	addressInformation.ApproximateLatitude = approximateLat

	approximateLong, err := offsetCoordinate(addressInformation.Longitude)
	if err != nil {
		return err
	}
	addressInformation.ApproximateLongitude = approximateLong

	osmAddress, err := nominatim.GetAddress(float32(approximateLat), float32(approximateLong))
	if err != nil {
		return err
	}

	addressInformation.ApproximateAddress = fmt.Sprintf("%s %s, %s %s", osmAddress.Road,
		osmAddress.HouseNumber, osmAddress.Postcode, osmAddress.City)
	addressInformation.ApproximateLatLong = fmt.Sprintf("%.8g;%.8g", approximateLat, approximateLong)

	if strings.TrimSpace(osmAddress.HouseNumber) == "" {
		return getApproximateAddress(addressInformation)
	}

	return nil
}

func offsetCoordinate(coordinate float64) (float64, error) {
	limit := big.NewInt(6)
	n := int64(0)

	for {
		nBig, err := rand.Int(rand.Reader, limit)
		if err != nil {
			return 0, err
		}

		n = nBig.Int64() - 3

		if n != 0 {
			break
		}
	}

	return coordinate + 0.001*float64(n), nil
}
