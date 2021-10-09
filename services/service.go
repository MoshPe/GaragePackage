package services

import (
	"bufio"
	"fmt"
	Resource "github.com/MoshPe/GaragePackage/resources"
	Utils "github.com/MoshPe/GaragePackage/utils"
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
	if len(Resource.GetResources()) == 0{
		fmt.Println("Resources are needed to be imported or created first!!")
		return
	}
	//TODO
	//getImportSelect := getImportSelection("services")
	getImportSelect := 1
	switch getImportSelect {
	case Utils.ImportViaTextFile:
		//TODO
		//importServicesViaTxt(getFileName())
		importServicesViaTxt(inputServices)
	case Utils.AddManually:
		var getService Service
		var getServiceId int
		for ok := true; ok; {
			Utils.IntInput("Please enter the service id ->: ", "Wrong input service id", &getServiceId)
			ok = IsServiceExist(getServiceId)
		}
		getService.Name = Utils.InputName("service")

		Utils.IntInput("Please enter the service work time in Hrs ->: ",
			"Wrong input service work time", &getService.TimeHr)

		Utils.IntInput("Please enter the amount of resources ->:",
			"Wrong input resource's quantity", &getService.AmountResourcesNeeded)

		var resourceId int
		fmt.Println("Please enter the resources id's ->:")
		for i := 0; i < getService.AmountResourcesNeeded; i++ {
			Utils.IntInput("", "couldn't input resource's id", &resourceId)
			if !Resource.IsResourceExist(resourceId) {
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
	defer Utils.CloseFile(importFile)
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
			if !Resource.IsResourceExist(resourceId) {
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
	if *getServiceId, _ = strconv.Atoi(resources[serviceId]); !Utils.IsIntPositive(*getServiceId) {
		errResult = "Invalid given service id!"
	}
	if IsServiceExist(*getServiceId) {
		errResult = "Invalid given service id!"
	}
	if getService.Name = resources[serviceName]; !Utils.IsProductNameValid(getService.Name) {
		errResult = "product name -" + resources[serviceName] + " need to contain only letters a-z , A-Z"
	}

	if getService.TimeHr, _ = strconv.Atoi(resources[serviceWorkTime]); !Utils.IsIntPositive(getService.TimeHr) {
		errResult = "Invalid given service work time! - service id: " + strconv.Itoa(*getServiceId)
	}
	if getService.AmountResourcesNeeded, _ = strconv.Atoi(resources[serviceResourceQuantity]); !Utils.IsIntPositive(getService.AmountResourcesNeeded) {
		errResult = "Invalid given service's resource quantity!" + strconv.Itoa(*getServiceId)
	}
	return
}

func IsServiceExist(serviceId int) bool {
	if _, ok := serviceList[serviceId]; ok {
		return true
	}
	return false
}

func GetServiceById(serviceId int) Service {
	return serviceList[serviceId]
}

func PrintServices(fileToPrint *bufio.Writer) {
	for id, service := range serviceList {
		fmt.Printf("ID: %d, service name: %s, service work time %d ,resource amount needed: %d, resources id's [",
			id, service.Name, service.TimeHr, service.AmountResourcesNeeded)
		printServiceResourceList(service,fileToPrint)
		fmt.Printf("\n")
	}
}

func printServiceResourceList(service Service, fileToPrint *bufio.Writer){
	for _, resourceId := range service.ResourcesIdList{
		fmt.Printf("%d ",resourceId)
	}
	fmt.Printf("]")
}

func PrintServiceNeededResources(serviceId int ,fileToPrint *bufio.Writer){
	service := serviceList[serviceId]
	fmt.Printf("service name %s resources->: [",service.Name)
	printServiceResourceList(service,fileToPrint)
	fmt.Printf("\n")
}
