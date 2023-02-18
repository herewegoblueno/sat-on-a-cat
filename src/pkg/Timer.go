package pkg

import (
	"time"
)

var isRunning bool = false
var startingTime time.Time

func StartTimer() {
	isRunning = true
	startingTime = time.Now()
}

func StopTimer() {
	isRunning = false
}

func GetElapsedNano() float64 {
	var elapsed int = int(time.Now().UnixNano() - startingTime.UnixNano())
	return float64(elapsed)
}
