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
	timeArrival      time.Time
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
		importServicesViaTxt(getFileName())
	case addManually:
		var getRequest Request
		var getRequestId int
		for ok := true; ok; {
			intInput("Please enter the car id ->: ", "Wrong input for car id", &getRequestId)
			ok = isServiceExist(getRequestId)
		}

		fmt.Printf("Please enter the car time of arrival (hh:mm) ->: ")
		if _, err := fmt.Scanln(&getRequest.timeArrival); err != nil {
			log.Fatalln("Wrong input arrival time")
		}

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
	defer func(importFile *os.File) {
		err := importFile.Close()
		if err != nil {
			log.Fatalln("Error in closing the import file")
		}
	}(importFile)
	scanner := bufio.NewScanner(importFile)
	for scanner.Scan() {
		var getService Service
		var getServiceId int
		resources := strings.Split(scanner.Text(), "\t")
		getService.name = resources[1]
		if !isProductNameValid(getService.name) {
			fmt.Println("product name -" + resources[1] + " need to contain only letters a-z , A-Z")
			continue
		}
		if getServiceId, _ = strconv.Atoi(resources[0]); isServiceExist(getServiceId) {
			fmt.Println("There is already a resource with the same id", getServiceId)
			continue
		}

		if !isIntPositive(getServiceId) {
			fmt.Println("Invalid given service id!")
			continue
		}
		getService.amountResourcesNeeded, _ = strconv.Atoi(resources[3])
		if !isIntPositive(getService.amountResourcesNeeded) {
			fmt.Println("Invalid given service's resources amount!")
			continue
		}
		for i := 0; i < getService.amountResourcesNeeded; i++ {
			serviceId, _ := strconv.Atoi(resources[i+4])
			getService.resourcesIdList = append(getService.resourcesIdList, serviceId)
		}
		serviceList[getServiceId] = getService
	}
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
