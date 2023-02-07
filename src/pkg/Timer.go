package pkg

import (
	"fmt"
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

func GetElapsedSeconds() (error, int) {
	if !isRunning {
		return fmt.Errorf("Timer hasn't been started yet!"), -1
	}

	var elapsed int = time.Now().Second() - startingTime.Second()
	return nil, elapsed
}
