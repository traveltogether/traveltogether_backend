package general

import "github.com/lib/pq"

func RemoveIntFromSlice(slice pq.Int64Array, object int64) (newSlice pq.Int64Array) {
	for _, element := range slice {
		if element != object {
			newSlice = append(newSlice, element)
		}
	}

	return
}
