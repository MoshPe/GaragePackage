package main

import (
	"fmt"
	Garage "github.com/MoshPe/GaragePackage"
	"sync"
	"time"
)

var (
	t time.Time
	resourcesScanLock = &sync.Mutex{}
	resourceLIst  map[int]Garage.Resource
	requestList  map[int]Garage.Request
)

func main(){
	fmt.Printf("\n\n---------- Welcome to our Garage, Enjoy :) ----------\n\n")
	initGarage()
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(len(requestList))
	for carId := range requestList{
		go waitForService(carId, waitGroup)
	}
	waitGroup.Wait()
	fmt.Printf("\n\n---------- ty :) ----------\n\n")

}

func waitForService(carId int, wg *sync.WaitGroup) {
	defer wg.Done()
	for t.Sub(requestList[carId].ArrivalTime) != 0 {}
	fmt.Println("car : ",carId," time : ",t.Format("15:04")," -> : arrived at the shop")
	fmt.Println("car : ",carId," time : ",t.Format("15:04")," -> : starting to collect resources")
	executeServices(carId)
}

func executeServices(carId int) {
	request := requestList[carId]
	for len(request.ServicesIdList) != 0 {
		for i,serviceId := range request.ServicesIdList{
			service := Garage.GetServiceById(serviceId)
			if scanResources(service){
				fmt.Println("car : ",carId," time : ",t.Format("15:04")," -> : started ",service.Name )
				time.Sleep(time.Second * time.Duration(service.TimeHr * 5))
				fmt.Println("car : ",carId," time : ",t.Format("15:04")," -> : finished ",service.Name )
				request.ServicesIdList = remove(request.ServicesIdList, i)
				//fmt.Println("car : ",carId," time : ",t.Format("15:04")," -> : ", request.ServicesIdList)
			}
		}
	}
}

func remove(slice []int, s int) []int {
	if len(slice) == 1 {
		return nil
	}
	if len(slice) {
		
	}
	return append(slice[:s], slice[s+1:]...)
}

func initGarage(){
	Garage.ImportResources()
	Garage.ImportServices()
	Garage.ImportRequests()
	resourceLIst = Garage.GetResources()
	requestList = Garage.GetRequests()
	go watcher()
	Garage.PrintResources()
	Garage.PrintServices()
	Garage.PrintRequests()
}

func scanResources(service Garage.Service)  bool{

	resourcesScanLock.Lock()
	{
			for _, resourceId := range service.ResourcesIdList {
				//means that we can't execute the request
				if resourceLIst[resourceId].AmountAvailable == 0 {
					resourcesScanLock.Unlock()
					return false
				}
			}
			takeResources(service)
	}
	resourcesScanLock.Unlock()
	return true
}

func takeResources(service Garage.Service) {
	for _, resourceId := range service.ResourcesIdList {
		resourceDown := resourceLIst[resourceId]
		resourceDown.AmountAvailable--
		resourceLIst[resourceId] = resourceDown
	}
}

func watcher() {
	//Initiating timer
	t ,_ = time.Parse("15:04","05:00")
	for true {
		time.Sleep(time.Second)
		t = t.Add(time.Minute * 15)
	}
}