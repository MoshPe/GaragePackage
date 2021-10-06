package Garage

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)
//TODO using sync.Map for concurrency

// The key is the car id
var requestList = sync.Map{}

func GetRequests() *sync.Map {
	return &requestList
}

//TODO just for making the testing easier at first
const inputRequests = "requests"


func ImportRequests() {
	if len(GetServices()) == 0 {
		fmt.Println("Services are needed to be imported or created first!!")
		return
	}
	//TODO
	//getImportSelect := getImportSelection("requests")
	getImportSelect := 1
	switch getImportSelect {
	case importViaTextFile:
		//TODO
		//importResourcesViaTxt(getFileName())
		importResourcesViaTxt(inputRequests)
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
		getRequest.ArrivalTime,_ = time.Parse("15:04",getArrivalTime)
		intInput("Please enter the amount of services ->:",
			"Wrong input service's quantity", &getRequest.AmountOfServices)

		var serviceId int
		fmt.Println("Please enter the services id's ->:")
		for i := 0; i < getRequest.AmountOfServices; i++ {
			intInput("", "couldn't input service's id", &serviceId)
			if !isServiceExist(serviceId) {
				fmt.Println("Service ", serviceId, " doesnt exist, Please try again!")
				i--
				continue
			}
			getRequest.ServicesIdList = append(getRequest.ServicesIdList, serviceId)
		}
		requestList.Store(getRequestId,getRequest)
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
		for i := 3; i < len(resources); i++ {
			serviceId, _ := strconv.Atoi(resources[i])
			if !isServiceExist(serviceId) {
				fmt.Println("Service ", serviceId, " doesnt exist, Please fix the file!. service id: ",getRequestId)
				getRequest.ServicesIdList = nil
				break
			}
			getRequest.ServicesIdList = append(getRequest.ServicesIdList, serviceId)
		}
		requestList.Store(getRequestId,getRequest)
		getRequest.ServicesIdList = nil
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
	if getRequest.ArrivalTime, err = time.Parse("15:04", resources[requestArrivalTime]); err != nil {
		errResult = "request arrival time -" + resources[requestArrivalTime] + " need to be as format hh:mm"
	}
	if getRequest.AmountOfServices, _ = strconv.Atoi(resources[requestResourceQuantity]); !isIntPositive(getRequest.AmountOfServices) {
		errResult = "Invalid given resource quantity!"
	}
	return
}

func isRequestExist(requestId int) bool {
	_, isExist := requestList.Load(requestId)
	return isExist
}


func PrintRequests() {
	requestList.Range(func(key, value interface{}) bool {
		requestId := key.(int)
		request := value.(Request)
		fmt.Printf("ID: %d, request Arrival Time name: %s, services amount needed: %d, services id's",requestId ,request.ArrivalTime.Format("15:04"),request.AmountOfServices)
		fmt.Println(request.ServicesIdList)
		return true
	})
}


/*

// The key is the car id
var requestList = make(map[int]Request)

func GetRequests() map[int]Request {
	return requestList
}

//TODO just for making the testing easier at first
const inputRequests = "requests"


func ImportRequests() {
	if len(GetServices()) == 0 {
		fmt.Println("Services are needed to be imported or created first!!")
		return
	}
	//TODO
	//getImportSelect := getImportSelection("requests")
	getImportSelect := 1
	switch getImportSelect {
	case importViaTextFile:
		//TODO
		//importResourcesViaTxt(getFileName())
		importResourcesViaTxt(inputRequests)
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
		getRequest.ArrivalTime,_ = time.Parse("15:04",getArrivalTime)
		intInput("Please enter the amount of services ->:",
			"Wrong input service's quantity", &getRequest.AmountOfServices)

		var serviceId int
		fmt.Println("Please enter the services id's ->:")
		for i := 0; i < getRequest.AmountOfServices; i++ {
			intInput("", "couldn't input service's id", &serviceId)
			if !isServiceExist(serviceId) {
				fmt.Println("Service ", serviceId, " doesnt exist, Please try again!")
				i--
				continue
			}
			getRequest.ServicesIdList = append(getRequest.ServicesIdList, serviceId)
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
		for i := 3; i < len(resources); i++ {
			serviceId, _ := strconv.Atoi(resources[i])
			if !isServiceExist(serviceId) {
				fmt.Println("Service ", serviceId, " doesnt exist, Please fix the file!. service id: ",getRequestId)
				getRequest.ServicesIdList = nil
				break
			}
			getRequest.ServicesIdList = append(getRequest.ServicesIdList, serviceId)
		}
		requestList[getRequestId] = getRequest
		getRequest.ServicesIdList = nil
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
	if getRequest.ArrivalTime, err = time.Parse("15:04", resources[requestArrivalTime]); err != nil {
		errResult = "request arrival time -" + resources[requestArrivalTime] + " need to be as format hh:mm"
	}
	if getRequest.AmountOfServices, _ = strconv.Atoi(resources[requestResourceQuantity]); !isIntPositive(getRequest.AmountOfServices) {
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


func PrintRequests() {
	for id,request := range requestList{
		fmt.Printf("ID: %d, request Arrival Time name: %s, services amount needed: %d, services id's",id,request.ArrivalTime.Format("15:04"),request.AmountOfServices)
		fmt.Println(request.ServicesIdList)
	}
}
 */