package resources

import (
	"time"
)

type Resource struct {
	Name            string
	AmountAvailable int
	WhenAvailable   [][]RequestTime
}

type RequestTime struct {
	CarId        int
	ServiceId    int
	FinishedTIme time.Time
}



