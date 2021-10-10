// Copyright Some Company Corp.
// All Rights Reserved// Here is where we explain the package.
// Some other stuff.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
	"unsafe"

	Request "github.com/MoshPe/GaragePackage/internal/requests"
	Resource "github.com/MoshPe/GaragePackage/internal/resources"
	Service "github.com/MoshPe/GaragePackage/internal/services"
	Utils "github.com/MoshPe/GaragePackage/pkg/utils"
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
	resourceLIst      map[int]Resource.Resource
	requestList       sync.Map
	printToFile       *bufio.Writer
)

func main() {
	defer func() {
		err := printToTxtFile().Close()
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
				// Resource.PrintResourcesShort(printToFile)
				// Service.PrintServiceNeededResources(serviceId,printToFile)
				printToLogs(carId, "starting to collect resources")
				changeResourceAmount(service, false)
				resourcesScanLock.Unlock()
				printToLogs(carId, "started "+service.Name)
				time.Sleep(time.Second * time.Duration(service.TimeHr*4))
				printToLogs(carId, "finished "+service.Name)
				carRequest.ServicesIdList = removeItem(carRequest.ServicesIdList, i, carId)
				resourcesScanLock.Lock()
				changeResourceAmount(service, true)
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
	for _, resourceID := range service.ResourcesIdList {
		// means that we can't execute the request
		if resourceLIst[resourceID].AmountAvailable == 0 {
			return false
		}
	}
	return true
}

func changeResourceAmount(service Service.Service, inc bool) {
	for _, resourceId := range service.ResourcesIdList {
		resourceDown := resourceLIst[resourceId]
		if inc {
			resourceDown.AmountAvailable++
		} else {
			resourceDown.AmountAvailable--
		}
		resourceLIst[resourceId] = resourceDown
	}
}

func initGarage() {
	resourceLIst = Resource.GetResources()
	requestList = *Request.GetRequests()
}

func startTimer() {
	//Initiating timer
	t, _ = time.Parse("15:04", "06:30")
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
	//fmt.Printf("\ncar : %d time : %s -> : %s",carId,t.Format("15:04"), msg)
	//fmt.Fprintf(printToFile,"\ncar : %d time : %s -> : %s",carId,t.Format("15:04"), msg)
	currentTime := C.CString(t.Format("15:04"))
	defer C.free(unsafe.Pointer(currentTime))
	msgToPrint := C.CString(msg)
	defer C.free(unsafe.Pointer(msgToPrint))
	C.printToLog(C.int(carId), (*C.char)(currentTime), (*C.char)(msgToPrint))
}

func printToTxtFile() *os.File {
	// TODO fix
	filePointer, err := os.Create("genericOutput.txt")

	if err != nil {
		log.Fatal(err)
	}
	printToFile = bufio.NewWriter(filePointer)
	//printToFile = bufio.NewWriter(os.Stdout)
	return filePointer
}

func endOfDay() bool {
	timesUp, _ := time.Parse("15:04", "15:30")
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
	//Resource.PrintResources(printToFile)
	//fmt.Println()
	//Service.PrintServices(printToFile)
	//fmt.Println()
	//Request.PrintRequests(printToFile)
	//fmt.Println()
	isRequestsEmpty := true
	exitProgram := false
	for isRequestsEmpty && !exitProgram {
		isRequestsEmpty = false
		C.printDayCountToLog(C.int(dayCount))
		//fmt.Fprintf(printToFile,"\n\n---------- Day %d :) ----------\n\n",dayCount)
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
