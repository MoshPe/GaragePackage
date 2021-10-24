// Package Garage Package Test Copyright Some Company Corp.
// All Rights Reserved// Here is where we explain the package.
// Some other stuff.
package Garage

import (
	"bufio"
	"fmt"
	Request "github.com/MoshPe/GaragePackage/internal/requests"
	Resource "github.com/MoshPe/GaragePackage/internal/resources"
	Service "github.com/MoshPe/GaragePackage/internal/services"
	Queue "github.com/MoshPe/GaragePackage/pkg/queue"
	Utils "github.com/MoshPe/GaragePackage/pkg/utils"
	"log"
	"os"
	"sync"
	"time"
	"unsafe"
)

/*
#include <stdio.h>
#include <stdlib.h>
#include "fileWrite.h"
*/
import "C"

var (
	t                 time.Time
	resourcesScanLock = &sync.Mutex{}
	resourceList      map[int]Resource.Resource
	requestList       sync.Map
	printToFile       *bufio.Writer
)

const(
	EndOfDayInTime = "15:30"
	beginningOfTheDay = "06:30"
)

func Run() {
	C.deleteFileIfExist()
	hmm := printToTxtFile()
	defer func() {
		err := hmm.Close()
		if err != nil {
			log.Fatal("Couldn't close the file")
		}
	}()
	fmt.Printf("\n\n---------- Welcome from our Garage, Enjoy :) ----------\n\n")

	menu()

	fmt.Printf("\n\n---------- ty :) ----------\n\n")
}

func waitForService(carId int, carRequest Request.Request, wg *sync.WaitGroup) {
	defer wg.Done()
	for t.Before(carRequest.ArrivalTime) {
	}
	printToLogs(carId, "arrived at the shop")
	executeServices(carId, carRequest)
}

func executeServices(carId int, carRequest Request.Request) {
	for len(carRequest.ServicesIdList) != 0 {
		for i, serviceId := range carRequest.ServicesIdList {
			service := Service.GetServiceById(serviceId)
			if endOfDay() {
				carRequest.ArrivalTime, _ = time.Parse("15:04", "06:30")
				requestList.Store(carId, carRequest)
				return
			}
			resourcesScanLock.Lock()
			if scanResources(service) {
				log.Printf("\ntimeNow %s - car %d - serviceId %d\n",t.Format("15:04"),carId,serviceId)
				Service.PrintServiceNeededResources(serviceId,printToFile)
				Resource.PrintResourcesShort(printToFile)
				printToLogs(carId, "starting to collect resources")
				changeResourceAmount(carId,serviceId,service, false, false)
				resourcesScanLock.Unlock()
				printToLogs(carId, "started "+service.Name)
				time.Sleep(time.Second * time.Duration(service.TimeHr*4))
				printToLogs(carId, "finished "+service.Name)
				carRequest.ServicesIdList = removeItem(carRequest.ServicesIdList, i, carId)
				resourcesScanLock.Lock()
				changeResourceAmount(carId,serviceId,service, true, false)
				resourcesScanLock.Unlock()
				break
			} else {
				resourcesScanLock.Unlock()
			}
		}
	}
	requestList.Delete(carId)
}

func scanResources(service Service.Service) bool {
	//TODO for not VIP requests we need to check for nil in every resource, check for having at least 1 empty timeQueue
	for _, resourceID := range service.ResourcesIdList {
		if resourceList[resourceID].AmountAvailable == 0 {
			return false
		}
	}
	return true
}

func changeResourceAmount(carId,serviceId int,service Service.Service, inc bool, isVIP bool) {
	for _, resourceId := range service.ResourcesIdList {
		resource := resourceList[resourceId]
		if inc {
			resource.AmountAvailable++
			pullRequestOutTime(carId,serviceId,resource)
		} else {
			finishTime := t.Add(time.Hour * time.Duration(service.TimeHr))
			insertRequestOutTime(carId,serviceId,finishTime, resource,isVIP)
			resource.AmountAvailable--
		}
		resourceList[resourceId] = resource
	}
}

func insertRequestOutTime(carId,ServiceId int, FinishedTIme time.Time,resource Resource.Resource, isVIP bool) {
	var insertTime = Resource.RequestTime{
		CarId:        carId,
		ServiceId:    ServiceId,
		FinishedTIme: FinishedTIme,
	}
	if isVIP {
			minInd := findMinimumFinishTime(resource.WhenAvailable)
			resource.WhenAvailable[minInd] = Queue.Enqueue(resource.WhenAvailable[minInd], insertTime)
		}else {
			for i := 0; i < len(resource.WhenAvailable); i++ {
				if resource.WhenAvailable[i] == nil {
					resource.WhenAvailable[i] = Queue.Enqueue(resource.WhenAvailable[i], insertTime)
				return
			}
		}
	}
}

func findMinimumFinishTime(arrayOfTimeList [][]Resource.RequestTime) int{
	minTime, _ := time.Parse("15:04", EndOfDayInTime)
	var minInd int
	for i := 0; i < len(arrayOfTimeList); i++ {
		if arrayOfTimeList[i] != nil {
			queueMinTime := Queue.Tail(arrayOfTimeList[i])
			if checkIfMin(minTime,queueMinTime){
				minTime = queueMinTime.FinishedTIme
				minInd = i
			}
		}else {
			//Means that the list in empty and the resource is free
			return i
		}
	}
	return minInd
}

func pullRequestOutTime(carId,ServiceId int, resource Resource.Resource)  {
	for i := 0; i < len(resource.WhenAvailable); i++ {
		if resource.WhenAvailable[i] != nil {
			if resource.WhenAvailable[i][0].CarId == carId &&
				resource.WhenAvailable[i][0].ServiceId == ServiceId{
					resource.WhenAvailable[i] = Queue.Dequeue(resource.WhenAvailable[i])
					if len(resource.WhenAvailable[i]) == 0{
						resource.WhenAvailable[i] = nil
					}
				}
			}
		}
}

func checkIfMin(minTime time.Time, isMin Resource.RequestTime) bool {
	if minTime.Before(isMin.FinishedTIme) {
		return true
	}
	return false
}

func printToTxtFile() *os.File {
	filePointer, err := os.Create("resourceTimeAvailability.txt")

	if err != nil {
		log.Fatal(err)
	}
	printToFile = bufio.NewWriter(filePointer)
	log.SetOutput(printToFile)
	return filePointer
}

func initGarage() {
	resourceList = Resource.GetResources()
	requestList = *Request.GetRequests()
}

func startTimer() {
	//Initiating timer
	t, _ = time.Parse("15:04", beginningOfTheDay)
	const timeSkip = 15
	for {
		time.Sleep(time.Second)
		t = t.Add(time.Minute * timeSkip)
	}
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

func printToLogs(carId int, msg string) {
	currentTime := C.CString(t.Format("15:04"))
	defer C.free(unsafe.Pointer(currentTime))
	msgToPrint := C.CString(msg)
	defer C.free(unsafe.Pointer(msgToPrint))
	C.printToLog(C.int(carId), (*C.char)(currentTime), (*C.char)(msgToPrint))
	log.Println("car : " +string(carId)+ " time : "+t.Format("15:04")+" -> : "+msg)
}

func endOfDay() bool {
	timesUp, _ := time.Parse("15:04", EndOfDayInTime)
	if t.After(timesUp) || t.Sub(timesUp) == 0 {
		return true
	}
	return false
}

func menu() {
	const (
		addedViaText = -1
		getResource  = 1
		getService   = 2
		getRequest   = 3
		runProgram   = 4
		exit         = 5
	)
	var (
		getImportSelect int
		dayCount        = 1
	)

	for {
		Utils.IntInput("1 - Import/add a resource\n"+
			"2 - Import/add service\n"+
			"3 - Import/add a request\n"+
			"4 - run the program\n"+
			"5 - exit the program\n"+
			"-> : ",
			"Wrong import selection input",
			&getImportSelect)

		switch getImportSelect {
		case getResource:
			Resource.ImportResources()
		case getService:
			Service.ImportServices()
		case getRequest:
			Request.ImportRequests()
		case runProgram:
			goto Run
		case exit:
			return
		default:
			fmt.Println("Wrong menu number, please try again!!")
		}
	}
Run:
	waitGroup := &sync.WaitGroup{}
	go startTimer()
	initGarage()
	isRequestsEmpty := true
	exitProgram := false
	for isRequestsEmpty && !exitProgram {
		isRequestsEmpty = false
		C.printDayCountToLog(C.int(dayCount))
		fmt.Fprintf(printToFile,"\n\n---------- Day %d :) ----------\n\n",dayCount)
		dayCount++
		t, _ = time.Parse("15:04", "06:30")
		requestList.Range(func(key, value interface{}) bool {
			isRequestsEmpty = true
			waitGroup.Add(1)
			go waitForService(key.(int), value.(Request.Request), waitGroup)
			return true
		})
		for {
			Utils.IntInput("1 - Import/add a request\n"+
				"2 - wait fr the end of the day\n"+
				"-> : ",
				"Wrong import selection input",
				&getImportSelect)
			//getImportSelect += 3
			if getImportSelect == 1 {
				newRequest, newRequestId := Request.ImportRequests()
				if newRequestId != addedViaText {
					waitGroup.Add(1)
					go waitForService(newRequestId, newRequest, waitGroup)
				}
			}
			if getImportSelect == 2 {
				fmt.Println("-------waiting for all services to finish---------")
				//exitProgram = true
				break
			}
		}
		waitGroup.Wait()
	}
}
