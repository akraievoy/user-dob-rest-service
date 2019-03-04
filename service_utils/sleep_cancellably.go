package service_utils

import (
	"context"
	"log"
	"time"
)

func SleepCancellably(ctx context.Context, timeout time.Duration) bool {
	timeoutAt := time.Now().Add(timeout)
	timeSlot := time.Millisecond * 200
	for time.Now().Before(timeoutAt) {
		select {
		case <-ctx.Done():
			log.Printf("sleep interrupted by context termination")
			return true
		default:
		}

		tillTimeout := timeoutAt.Sub(time.Now())
		if tillTimeout <= 0 {
			break
		} else if tillTimeout < timeSlot {
			time.Sleep(tillTimeout)
		} else {
			time.Sleep(timeSlot)
		}
	}
	return false
}
