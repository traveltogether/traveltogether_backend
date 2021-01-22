package general

func RemoveIntFromSlice(slice []int, object int) (newSlice []int) {
	for _, element := range slice {
		if element != object {
			newSlice = append(newSlice, element)
		}
	}

	return
}
