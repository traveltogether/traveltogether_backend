package journey

import (
	"errors"
	"fmt"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/database"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
	"strconv"
	"time"
)

var (
	NotFound                          = errors.New("journey not found")
	DeletionNotAvailableDueToRequests = errors.New("deletion not available due to requests")
)

func InsertJourneyToDatabase(journey *types.Journey) error {
	id := &types.IdInformation{}

	err := database.NamedQueryAsync(database.DefaultTimeout, id, ""+
		"INSERT INTO "+
		"journeys(user_id, request, offer, start_lat_long, end_lat_long, approximate_start_lat_long,"+
		"approximate_end_lat_long, start_address, end_address, approximate_start_address, approximate_end_address,"+
		"time_value, time_is_departure, time_is_arrival, open_for_requests, pending_user_ids, accepted_user_ids,"+
		"declined_user_ids, cancelled_by_host, cancelled_by_attendee_ids, note) "+
		"VALUES("+
		":user_id, :request, :offer, :start_lat_long, :end_lat_long, :approximate_start_lat_long,"+
		":approximate_end_lat_long, :start_address, :end_address, :approximate_start_address,"+
		":approximate_end_address, :time_value, :time_is_departure, :time_is_arrival, :open_for_requests,"+
		":pending_user_ids, :accepted_user_ids, :declined_user_ids, :cancelled_by_host, :cancelled_by_attendee_ids,"+
		":note) "+
		"RETURNING id",
		journey)

	if err != nil {
		return err
	}

	newId := &id.Id
	journey.Id = newId

	return nil
}

func RetrieveJourneyFromDatabase(id int) (*types.Journey, error) {
	slice, err := database.QueryAsync(database.DefaultTimeout, types.JourneyType,
		"SELECT * FROM journeys WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	journeys := slice.([]*types.Journey)
	if len(journeys) == 0 {
		return nil, NotFound
	}

	return journeys[0], nil
}

func GetAllJourneysFromDatabase() ([]*types.Journey, error) {
	slice, err := database.QueryAsync(database.DefaultTimeout, types.JourneyType, "SELECT * from journeys")
	if err != nil {
		return nil, err
	}

	return slice.([]*types.Journey), nil
}

func GetJourneysFromDatabaseFilteredBy(filters map[string]interface{}, nonExpiredFilter bool) ([]*types.Journey, error) {
	query := "SELECT * FROM journeys"
	filterLength := len(filters)
	values := make([]interface{}, 0, filterLength)
	if filterLength != 0 {
		query += " WHERE "
		index := 1
		for key, value := range filters {
			query += fmt.Sprintf("%s=$%d AND ", key, index)
			values = append(values, value)
			index += 1
		}
		query = query[:len(query)-5]
	}
	if nonExpiredFilter {
		if filterLength == 0 {
			query += " WHERE "
		} else {
			query += " AND "
		}
		query += " time_value > " + strconv.Itoa(int(time.Now().UnixNano()/int64(time.Millisecond)))
	}

	slice, err := database.QueryAsync(database.DefaultTimeout, types.JourneyType, query, values...)
	if err != nil {
		return nil, err
	}

	return slice.([]*types.Journey), nil
}

func DeleteJourneyFromDatabase(id int) error {
	slice, err := database.QueryAsync(database.DefaultTimeout, types.IdInformationType,
		"DELETE FROM journeys WHERE id = $1 "+
			"AND (pending_user_ids IS NULL OR pending_user_ids = '{}') "+
			"AND (accepted_user_ids IS NULL OR accepted_user_ids = '{}') "+
			"RETURNING id", id)

	if err != nil {
		return err
	}

	if len(slice.([]*types.IdInformation)) == 0 {
		return DeletionNotAvailableDueToRequests
	}

	return nil
}
