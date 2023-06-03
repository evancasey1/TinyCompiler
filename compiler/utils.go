package compiler

func setAdd(slice []string, elem string) []string {
	elemPresent := false
	for _, e := range slice {
		if e == elem {
			elemPresent = true
		}
	}

	if !elemPresent {
		slice = append(slice, elem)
	}
	return slice
}

func elementInSet(slice []string, elem string) bool {
	for _, e := range slice {
		if e == elem {
			return true
		}
	}
	return false
}
