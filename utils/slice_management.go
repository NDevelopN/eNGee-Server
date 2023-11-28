package utils

import "fmt"

func RemoveElementFromSlice[T comparable](slice []T, element T) ([]T, error) {
	for i, sliceEle := range slice {
		if element == sliceEle {
			slice[i] = slice[len(slice)-1]
			if len(slice) == 1 {
				slice = make([]T, 0)
			}
			return slice, nil
		}
	}

	return slice, fmt.Errorf("matching element not found in slice")
}

func RemoveElementFromSliceOrdered[T comparable](slice []T, element T) ([]T, error) {
	for i, sliceEle := range slice {
		if element == sliceEle {
			slice = append(slice[:i], slice[i+1:]...)
			return slice, nil
		}
	}

	return slice, fmt.Errorf("matching element not found in slice")
}
