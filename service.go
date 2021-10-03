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

func GetServices() *map[int]Service {
	return &serviceList
}

func ImportServices() {
	if GetResources() == nil {
		log.Fatalln("Resources are needed to be imported or created first!!")
		return
	}
	getImportSelect := getImportSelection("services")
	switch getImportSelect {
	case importViaTextFile:
		importServicesViaTxt(getFileName())
	case addManually:
		var getService Service
		var getServiceId int
		for ok := true; ok; {
			intInput("Please enter the service id ->: ", "Wrong input service id", &getServiceId)
			ok = isServiceExist(getServiceId)
		}
		getService.name = inputName("service")

		intInput("Please enter the service work time in Hrs ->: ",
			"Wrong input service work time", &getService.timeHr)

		intInput("Please enter the amount of resources ->:",
			"Wrong input resource's quantity", &getService.amountResourcesNeeded)

		var resourceId int
		fmt.Println("Please enter the resources id's ->:")
		for i := 0; i < getService.amountResourcesNeeded; i++ {
			intInput("", "couldn't input resource's id", &resourceId)
			if !isResourceExist(resourceId) {
				fmt.Println("Resource ", resourceId, " doesnt exist, Please try again!")
				i--
				continue
			}
			getService.resourcesIdList = append(getService.resourcesIdList, resourceId)
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
		for i := 0; i < getService.amountResourcesNeeded; i++ {
			resourceId, _ := strconv.Atoi(resources[i+4])
			if !isResourceExist(resourceId) {
				fmt.Println("Resource ", resourceId, " doesnt exist, Please fix the file!. service id: ",getServiceId)
				getService.resourcesIdList = nil
				break
			}
			getService.resourcesIdList = append(getService.resourcesIdList, resourceId)
		}

		//Update the map with the new service
		serviceList[getServiceId] = getService
		getService.resourcesIdList = nil
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
	if getService.name = resources[serviceName]; !isProductNameValid(getService.name) {
		errResult = "product name -" + resources[serviceName] + " need to contain only letters a-z , A-Z"
	}

	if getService.timeHr, _ = strconv.Atoi(resources[serviceWorkTime]); !isIntPositive(getService.timeHr) {
		errResult = "Invalid given service work time! - service id: " + string(*getServiceId)
	}
	if getService.amountResourcesNeeded, _ = strconv.Atoi(resources[serviceResourceQuantity]); !isIntPositive(getService.amountResourcesNeeded) {
		errResult = "Invalid given service's resource quantity!" + string(*getServiceId)
	}
	return
}

func isServiceExist(serviceId int) bool {
	if _, ok := serviceList[serviceId]; ok {
		return true
	}
	return false
}

func getServiceById(serviceId int) Service {
	return serviceList[serviceId]
}

func PrintServices() {
	for id, service := range serviceList {
		fmt.Printf("ID: %d, service name: %s, service work time %d ,resource amount needed: %d, resources id's", id, service.name, service.timeHr, service.amountResourcesNeeded)
		fmt.Println(service.resourcesIdList)
	}
}
