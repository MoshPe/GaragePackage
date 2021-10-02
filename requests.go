package Garage

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Request struct {
	arrivalTime      time.Time
	amountOfServices int
	servicesId       []int
}
// The key is the car id
var requestList = make(map[int]Request)

func GetRequests() *map[int]Request {
	return &requestList
}

func ImportRequests() {
	getImportSelect := getImportSelection("requests")
	switch getImportSelect {
	case importViaTextFile:
		importRequestsViaTxt(getFileName())
	case addManually:
		var getRequest Request
		var getRequestId int
		var getArrivalTime string
		for ok := true; ok; {
			intInput("Please enter the car id ->: ", "Wrong input for car id", &getRequestId)
			ok = isServiceExist(getRequestId)
		}

		fmt.Printf("Please enter the car time of arrival (hh:mm) ->: ")
		if _, err := fmt.Scanln(&getArrivalTime); err != nil {
			log.Fatalln("Wrong input arrival time")
		}
		getRequest.arrivalTime,_ = time.Parse("15:04",getArrivalTime)
		intInput("Please enter the amount of services ->:",
			"Wrong input service's quantity", &getRequest.amountOfServices)

		var serviceId int
		fmt.Println("Please enter the services id's ->:")
		for i := 0; i < getRequest.amountOfServices; i++ {
			intInput("", "couldn't input service's id", &serviceId)
			if !isServiceExist(serviceId) {
				fmt.Println("Service ", serviceId, " doesnt exist, Please try again!")
				i--
				continue
			}
			getRequest.servicesId = append(getRequest.servicesId, serviceId)
		}
		requestList[getRequestId] = getRequest
	}
}

func importRequestsViaTxt(fileName string) {
	importFile, err := os.Open(fileName + ".txt") // For read access.
	if err != nil {
		log.Fatal(err)
	}
	//close the file when the function finishes
	defer closeFile(importFile)
	scanner := bufio.NewScanner(importFile)
	var getRequest Request
	var getRequestId int
	for scanner.Scan() {
		resources := strings.Split(scanner.Text(), "\t")
		if errResult := checkRequestValidation(resources, &getRequestId, &getRequest); errResult != "" {
			fmt.Println(errResult)
		}
		for i := 0; i < getRequest.amountOfServices; i++ {
			serviceId, _ := strconv.Atoi(resources[i+3])
			getRequest.servicesId = append(getRequest.servicesId, serviceId)
		}
		requestList[getRequestId] = getRequest
	}
}

func checkRequestValidation(resources []string, getRequestId *int, getRequest *Request) (errResult string) {
	const (
		requestId               = 0
		requestArrivalTime      = 1
		requestResourceQuantity = 2
	)
	var err error
	if *getRequestId, _ = strconv.Atoi(resources[requestId]); !isIntPositive(*getRequestId) {
		errResult = "Invalid given request id!"
	}
	if isRequestExist(*getRequestId) {
		errResult = "Invalid given resource id!"
	}
	if getRequest.arrivalTime, err = time.Parse("15:04", resources[requestArrivalTime]); err != nil {
		errResult = "request arrival time -" + resources[requestArrivalTime] + " need to be as format hh:mm"
	}
	if getRequest.amountOfServices, _ = strconv.Atoi(resources[requestResourceQuantity]); !isIntPositive(getRequest.amountOfServices) {
		errResult = "Invalid given resource quantity!"
	}
	return
}

func isRequestExist(requestId int) bool {
	if _, ok := requestList[requestId]; ok {
		return true
	}
	return false
}

/*
func PrintServices() {
	for id,service := range serviceList{
		fmt.Printf("ID: %d, service name: %s, service work time %f ,resource amount needed: %d, resources id's",id,service.name,service.timeHr,service.amountResourcesNeeded)
		fmt.Println(service.resourcesIdList)
	}
}
*/
