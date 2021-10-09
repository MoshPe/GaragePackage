package requests

import "time"

type Request struct {
	ArrivalTime time.Time
	AmountOfServices int
	ServicesIdList []int
}