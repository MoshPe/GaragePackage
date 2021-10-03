package main

import (
	"fmt"
	Garage "github.com/MoshPe/GaragePackage"
	"os"
	"sync"
	"time"
)

var (
	t time.Time
	resourcesScanLock = &sync.Mutex{}
	resourceAmountChange = &sync.Mutex{}
	requestRW = &sync.Mutex{}
	resourceLIst  map[int]Garage.Resource
	requestList  map[int]Garage.Request
)

func main(){
	defer Garage.PrintRequests()
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
	requestRW.Lock()
	request := requestList[carId]
	requestRW.Unlock()
	for len(request.ServicesIdList) != 0 {
		for i,serviceId := range request.ServicesIdList{
			service := Garage.GetServiceById(serviceId)
			resourcesScanLock.Lock()
			if scanResources(service){
				changeResourceAmount(service,false)
				resourcesScanLock.Unlock()
				fmt.Println("car : ",carId," time : ",t.Format("15:04")," -> : started ",service.Name )
				time.Sleep(time.Millisecond * time.Duration(service.TimeHr * 10))
				fmt.Println("car : ",carId," time : ",t.Format("15:04")," -> : finished ",service.Name )
				request.ServicesIdList = removeItem(request.ServicesIdList, i, carId)
				//fmt.Println("car : ",carId," time : ",t.Format("15:04")," -> : ", request.ServicesIdList)
				resourcesScanLock.Lock()
				changeResourceAmount(service,true)
				resourcesScanLock.Unlock()
				break
			}else{
				resourcesScanLock.Unlock()
			}
		}
	}
	requestRW.Lock()
	requestList[carId] = request
	requestRW.Unlock()
}

func removeItem(slice []int, s int, carId int) []int {
	defer func() { //catch or finally
		if err := recover(); err != nil { //catch
			fmt.Fprintf(os.Stderr, "Exception: %d\n", carId)
			os.Exit(1)
		}
	}()

	if len(slice) == 1 {
		return nil
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
	fmt.Println()
	Garage.PrintServices()
	fmt.Println()
	Garage.PrintRequests()
	fmt.Println()
}

func scanResources(service Garage.Service)  bool{
	{
		for _, resourceId := range service.ResourcesIdList {
			//means that we can't execute the request
			if resourceLIst[resourceId].AmountAvailable == 0 {
				return false
			}
		}
	}
	return true
}

func changeResourceAmount(service Garage.Service, inc bool) {
	for _, resourceId := range service.ResourcesIdList {
		resourceDown := resourceLIst[resourceId]
		if inc {
			resourceDown.AmountAvailable++
		}else {
			resourceDown.AmountAvailable--
		}
		resourceLIst[resourceId] = resourceDown
	}
}

func watcher() {
	//Initiating timer
	t ,_ = time.Parse("15:04","06:30")
	for true {
		time.Sleep(time.Second)
		t = t.Add(time.Minute * 15)
	}
}