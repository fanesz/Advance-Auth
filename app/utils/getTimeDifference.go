package utils

import (
	"time"
)

func GetTimeDifference(inputDate time.Time) int {
	currentTime := time.Now()
	difference := currentTime.Sub(inputDate).Minutes()

	return int(difference)
}
