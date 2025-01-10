package expire

import (
	"errors"
	"time"
)

// TimeDate represents a Unix timestamp as an expiration date
type TimeDate int64

// ExpireIn creates a new TimeDate that is `addMinutes` minutes in the future
func ExpireIn(addMinutes int) TimeDate {
	add := time.Duration(addMinutes) * time.Minute
	return TimeDate(time.Now().Add(add).Unix())
}

// IsExpired checks if the current TimeDate is in the past
func (d TimeDate) IsExpired() bool {
	return TimeDate(time.Now().Unix()) > d
}

// RemainingTime returns the time duration remaining until expiration
func (d TimeDate) RemainingTime() (time.Duration, error) {
	now := time.Now().Unix()
	if d.IsExpired() {
		return 0, errors.New("time has already expired")
	}
	return time.Duration(d-TimeDate(now)) * time.Second, nil
}

// ExpireAt creates a TimeDate object for a specific time
func ExpireAt(target time.Time) TimeDate {
	return TimeDate(target.Unix())
}
