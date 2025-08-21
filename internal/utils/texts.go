package utils

import (
	"fmt"
	"slices"
)

func GetWateringInterval(wateringInterval int) string {
	switch {
	case wateringInterval%10 == 1:
		return fmt.Sprintf("%d день", wateringInterval)
	case slices.Contains([]int{2, 3, 4}, wateringInterval):
		return fmt.Sprintf("%d дня", wateringInterval)
	default:
		return fmt.Sprintf("%d дней", wateringInterval)
	}
}
