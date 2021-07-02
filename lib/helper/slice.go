package helper

func IsInSlice(slice []string, val string) (int, bool) {
	if len(slice) == 0 {
		return -1, false
	}
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
