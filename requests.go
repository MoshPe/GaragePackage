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
	timeArrival time.Time
	amountOfServices int
	servicesId []int
}

// The key is the car id
var requestList = make(map[int] Request)

func GetRequests() *map[int] Request{
	return &requestList
}

func ImportRequests(){
	getImportSelect := getImportSelection("requests")
	switch getImportSelect {
	case importViaTextFile:
		importServicesViaTxt(getFileName())
	case addResourceManually:
		var getService Service
		var getServiceId int
		for ok := true; ok ;{
			intInput("Please enter the service id ->: ","Wrong input service id",&getServiceId)
			ok = isServiceExist(getServiceId)
		}
		fmt.Printf("Please enter the service work time in Hrs ->: ")
		if _,err := fmt.Scanln(&getService.timeHr); err != nil {
			log.Fatalln("Wrong input service work time")
		}
		fmt.Println("Please enter the amount of resources ->:")
		if _,err := fmt.Scanln(&getService.amountResourcesNeeded); err != nil {
			log.Fatalln("Wrong input resource's quantity")
		}
		var resourceId int
		fmt.Println("Please enter the resources id's ->:")
		for i := 0;i < getService.amountResourcesNeeded;i++ {
			if _,err := fmt.Scanf("%d",&resourceId); err != nil{
				log.Fatalln("couldn't input resource's id")
			}
			if !isResourceExist(resourceId){
				fmt.Println("Resource",resourceId,"doesnt exist, Please try again")
				i--
				continue
			}
			getService.resourcesIdList = append(getService.resourcesIdList,resourceId)
		}
		serviceList[getServiceId] = getService
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
	for scanner.Scan(){
		var getService Service
		var getServiceId int
		resources := strings.Split(scanner.Text(), "\t")
		getService.name = resources[1]
		if !isProductNameValid(getService.name) {
			fmt.Println("product name -"+ resources[1] +" need to contain only letters a-z , A-Z")
			continue
		}
		if getServiceId,_ = strconv.Atoi(resources[0]); isServiceExist(getServiceId){
			fmt.Println("There is already a resource with the same id", getServiceId)
			continue
		}

		if !isIntPositive(getServiceId) {
			fmt.Println("Invalid given service id!")
			continue
		}
		getService.timeHr, _ =strconv.ParseFloat(resources[2],64)
		getService.amountResourcesNeeded, _ =strconv.Atoi(resources[3])
		if !isIntPositive(getService.amountResourcesNeeded) {
			fmt.Println("Invalid given service's resources amount!")
			continue
		}
		for i := 0; i < getService.amountResourcesNeeded; i++ {
			serviceId,_ := strconv.Atoi(resources[i+4])
			getService.resourcesIdList = append(getService.resourcesIdList,serviceId)
		}
		serviceList[getServiceId] = getService
	}
}

func isRequestExist(requestId int) bool{
	if _, ok := requestList[requestId]; ok{
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
