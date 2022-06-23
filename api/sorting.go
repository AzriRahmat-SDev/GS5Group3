package api

import "strings"

func sortVenueInfoList(v []VenueInformation, left int, right int) {
	if left < right {
		pivotIndex := venueInfoPartition(v, left, right)
		sortVenueInfoList(v, left, pivotIndex-1)
		sortVenueInfoList(v, pivotIndex+1, right)
	}
}

func venueInfoPartition(v []VenueInformation, left int, right int) int {
	pivot := strings.ToLower(v[(left+right)/2].VenueName)
	for left < right {
		for strings.ToLower(v[left].VenueName) < pivot {
			left++
		}
		for strings.ToLower(v[right].VenueName) > pivot {
			right--
		}
		if left < right {
			temp := v[left]
			v[left] = v[right]
			v[right] = temp
		}
	}
	return left
}
