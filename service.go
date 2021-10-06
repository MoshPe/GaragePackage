package Garage

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var serviceList = make(map[int]Service)

//TODO just for making the testing easier at first
const inputServices = "services"


func GetServices() map[int]Service {
	return serviceList
}

func ImportServices() {
	if len(GetResources()) == 0{
		fmt.Println("Resources are needed to be imported or created first!!")
		return
	}
	//TODO
	//getImportSelect := getImportSelection("services")
	getImportSelect := 1
	switch getImportSelect {
	case importViaTextFile:
		//TODO
		//importServicesViaTxt(getFileName())
		importServicesViaTxt(inputServices)
	case addManually:
		var getService Service
		var getServiceId int
		for ok := true; ok; {
			intInput("Please enter the service id ->: ", "Wrong input service id", &getServiceId)
			ok = isServiceExist(getServiceId)
		}
		getService.Name = inputName("service")

		intInput("Please enter the service work time in Hrs ->: ",
			"Wrong input service work time", &getService.TimeHr)

		intInput("Please enter the amount of resources ->:",
			"Wrong input resource's quantity", &getService.AmountResourcesNeeded)

		var resourceId int
		fmt.Println("Please enter the resources id's ->:")
		for i := 0; i < getService.AmountResourcesNeeded; i++ {
			intInput("", "couldn't input resource's id", &resourceId)
			if !isResourceExist(resourceId) {
				fmt.Println("Resource ", resourceId, " doesnt exist, Please try again!")
				i--
				continue
			}
			getService.ResourcesIdList = append(getService.ResourcesIdList, resourceId)
		}
		serviceList[getServiceId] = getService
	}
}

func importServicesViaTxt(fileName string) {
	importFile, err := os.Open(fileName + ".txt") // For read access.
	if err != nil {
		log.Fatal(err)
	}

	//close the file when the function finishes
	defer closeFile(importFile)
	scanner := bufio.NewScanner(importFile)
	var getService Service
	var getServiceId int
	const serviceWorkTime = 2
	for scanner.Scan() {
		resources := strings.Split(scanner.Text(), "\t")
		if errResult := checkServiceValidation(resources, &getServiceId, &getService); errResult != "" {
			fmt.Println(errResult)
		}

		//Input the service's resources as a list
		for i := 0; i < getService.AmountResourcesNeeded; i++ {
			resourceId, _ := strconv.Atoi(resources[i+4])
			if !isResourceExist(resourceId) {
				fmt.Println("Resource ", resourceId, " doesnt exist, Please fix the file!. service id: ",getServiceId)
				getService.ResourcesIdList = nil
				break
			}
			getService.ResourcesIdList = append(getService.ResourcesIdList, resourceId)
		}

		//Update the map with the new service
		serviceList[getServiceId] = getService
		getService.ResourcesIdList = nil
	}
}

func checkServiceValidation(resources []string, getServiceId *int, getService *Service) (errResult string) {
	const (
		serviceId               = 0
		serviceName             = 1
		serviceWorkTime         = 2
		serviceResourceQuantity = 3
	)
	if *getServiceId, _ = strconv.Atoi(resources[serviceId]); !isIntPositive(*getServiceId) {
		errResult = "Invalid given service id!"
	}
	if isServiceExist(*getServiceId) {
		errResult = "Invalid given resource id!"
	}
	if getService.Name = resources[serviceName]; !isProductNameValid(getService.Name) {
		errResult = "product name -" + resources[serviceName] + " need to contain only letters a-z , A-Z"
	}

	if getService.TimeHr, _ = strconv.Atoi(resources[serviceWorkTime]); !isIntPositive(getService.TimeHr) {
		errResult = "Invalid given service work time! - service id: " + strconv.Itoa(*getServiceId)
	}
	if getService.AmountResourcesNeeded, _ = strconv.Atoi(resources[serviceResourceQuantity]); !isIntPositive(getService.AmountResourcesNeeded) {
		errResult = "Invalid given service's resource quantity!" + strconv.Itoa(*getServiceId)
	}
	return
}

func isServiceExist(serviceId int) bool {
	if _, ok := serviceList[serviceId]; ok {
		return true
	}
	return false
}

func GetServiceById(serviceId int) Service {
	return serviceList[serviceId]
}

func PrintServices() {
	for id, service := range serviceList {
		fmt.Printf("ID: %d, service name: %s, service work time %d ,resource amount needed: %d, resources id's", id, service.Name, service.TimeHr, service.AmountResourcesNeeded)
		fmt.Println(service.ResourcesIdList)
	}
}
