package Garage
/*
import (
	Request "github.com/MoshPe/GaragePackage/internal/requests"
	Resource "github.com/MoshPe/GaragePackage/internal/resources"
	Service "github.com/MoshPe/GaragePackage/internal/services"
	Queue "github.com/MoshPe/GaragePackage/pkg/queue"
	"time"
)

type VIPrequest struct {
	request Request.Request
	endRequestInTime time.Time
}

func predictConstrainSuccess(vip VIPrequest,resources map[int]Resource.Resource) bool {
	// That channel will hold all the services end run time
	requestStartTime := vip.request.ArrivalTime
	ch := make(chan time.Time, vip.request.AmountOfServices)
	serviceTimes := make([]time.Time,vip.request.AmountOfServices)
	serviceIds := vip.request.ServicesIdList
	for i:= 0;i < vip.request.AmountOfServices; i++{
		for i:= 0;i < vip.request.AmountOfServices; i++{
			service := Service.GetServiceById(serviceIds[i])
			go calculateEndTime(ch, vip,service,resources)
			for i:= 0;i < vip.request.AmountOfServices; i++ {
				serviceTimes[i] = <- ch
			}
			minServiceStartTime := findMinTime(serviceTimes)
			requestStartTime = requestStartTime.Add(time.Duration(time.Hour * time.Duration(minServiceStartTime.Hour()) +
																	time.Minute * time.Duration(minServiceStartTime.Minute())))
			//TODO need to remove the service id
			//removeItem(serviceIds,)
		}
	}

	return false
}

func calculateEndTime(ch chan time.Time,vip VIPrequest, service Service.Service, resources map[int]Resource.Resource) time.Time{
	//TODO if service.amountOfResources == 0 return
	if service.AmountResourcesNeeded == 0 {
		return vip.request.ArrivalTime.Add(time.Hour * time.Duration(service.TimeHr))
	}
	serviceStartTime := getServiceStartTime(service,resources)
	ch <- serviceStartTime.Add(time.Hour * time.Duration(service.TimeHr))
	return time.Time{}
}

func getServiceStartTime(service Service.Service, resources map[int]Resource.Resource) time.Time{
	var resourcesStartTime []time.Time
	for _, resourceId := range service.ResourcesIdList {
		minTime, _ := time.Parse("15:04", EndOfDayInTime)
		resourcesQueueTime := resources[resourceId].WhenAvailable
		for i := 0; i < len(resourcesQueueTime); i++ {
			if resourcesQueueTime[i] != nil {
				queueMinTime := Queue.Tail(resourcesQueueTime[i])
				if !checkIfMin(minTime,queueMinTime){
					minTime = queueMinTime.FinishedTIme
				}
			}
		}
		resourcesStartTime = append(resourcesStartTime, minTime)
	}
	return findMaxTime(resourcesStartTime)
}

func findMaxTime(times []time.Time) time.Time{
	maxTime := times[0]
	for _, t := range times {
		if maxTime.Before(t) {
			maxTime = t
		}
	}
	return maxTime
}

func findMinTime(times []time.Time) time.Time{
	maxTime := times[0]
	for _, t := range times {
		if maxTime.After(t) {
			maxTime = t
		}
	}
	return maxTime
}

*/