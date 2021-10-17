package queue

import "github.com/MoshPe/GaragePackage/internal/resources"

func Enqueue(queue[] resources.RequestTime, element resources.RequestTime) [] resources.RequestTime {
	queue = append(queue, element) // Simply append to enqueue.
	return queue
}

func Dequeue(queue[] resources.RequestTime) [] resources.RequestTime {
	if len(queue) == 0 {
		return nil
	}
	return  queue[1:]
}

func Head(queue[] resources.RequestTime) resources.RequestTime{
	if len(queue) == 0 {
		return resources.RequestTime{}
	}
	return queue[0]
}

func Tail(queue[] resources.RequestTime) resources.RequestTime{
	if len(queue) == 0 {
		return resources.RequestTime{}
	}
	return queue[len(queue) - 1]
}
