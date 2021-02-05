package journey

import (
	"errors"
	"fmt"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/database"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
)

var (
	NotFound = errors.New("journey not found")
)

func InsertJourneyToDatabase(journey *types.Journey) error {
	id := &types.IdInformation{}

	err := database.NamedQueryAsync(database.DefaultTimeout, id, ""+
		"INSERT INTO "+
		"journeys(user_id, request, offer, start_lat_long, end_lat_long, approximate_start_lat_long,"+
		"approximate_end_lat_long, start_address, end_address, approximate_start_address, approximate_end_address,"+
		"time_value, time_is_departure, time_is_arrival, open_for_requests, pending_user_ids, accepted_user_ids,"+
		"declined_user_ids, cancelled_by_host, cancelled_by_attendee_ids) "+
		"VALUES("+
		":user_id, :request, :offer, :start_lat_long, :end_lat_long, :approximate_start_lat_long,"+
		":approximate_end_lat_long, :start_address, :end_address, :approximate_start_address,"+
		":approximate_end_address, :time_value, :time_is_departure, :time_is_arrival, :open_for_requests,"+
		":pending_user_ids, :accepted_user_ids, :declined_user_ids, :cancelled_by_host, :cancelled_by_attendee_ids) "+
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

func GetJourneysFromDatabaseFilteredBy(filters map[string]interface{}) ([]*types.Journey, error) {
	query := "SELECT * FROM journeys"
	values := make([]interface{}, 0, len(filters))
	if len(filters) != 0 {
		query += " WHERE "
		index := 1
		for key, value := range filters {
			query += fmt.Sprintf("%s=$%d AND ", key, index)
			values = append(values, value)
			index += 1
		}
		query = query[:len(query)-5]
	}

	slice, err := database.QueryAsync(database.DefaultTimeout, types.JourneyType, query, values...)
	if err != nil {
		return nil, err
	}

	return slice.([]*types.Journey), nil
}

func DeleteJourneyFromDatabase(id int) error {
	return database.PrepareAsync(database.DefaultTimeout, "DELETE FROM journeys WHERE id = $1", id)
}
