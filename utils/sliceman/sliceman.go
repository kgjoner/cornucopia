package sliceman

// Return index number of target. If target is not in slice, it returns -1. 
func IndexOf[T comparable](slice []T, target T) int {
	res := -1
	for index, el := range slice {
		if el == target {
			res = index
		}
	}

	return res
}

// Remove element at index. It is faster than SafeRemove, but does not keep elements order.
func Remove[T any](slice []T, index int) []T {
	if index < 0 {
		return slice
	}

	slice[index] = slice[len(slice)-1]
	slice = slice[:len(slice)-1]
	return slice
}

// Remove element at index. Assure initial order is preserved.
func SafeRemove[T any](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}
