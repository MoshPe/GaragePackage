package resources

import (
	"bufio"
	"fmt"
	Utils "github.com/MoshPe/GaragePackage/pkg/utils"
	"log"
	"os"
	"strconv"
	"strings"
)

var resourcesList = make(map[int]Resource)

//TODO just for making the testing easier at first
const inputResources = "resources"

func GetResources() map[int]Resource {
	return resourcesList
}

func ImportResources() {
	//TODO
	//getImportSelect := getImportSelection("resources")
	getImportSelect := 1
	switch getImportSelect {
	case Utils.ImportViaTextFile:
		//TODO
		//importResourcesViaTxt(getFileName())
		importResourcesViaTxt(inputResources)
	case Utils.AddManually:
		var getResource Resource
		var getResourceId int
		for ok := true; ok; {
			Utils.IntInput("Please enter the resource id ->: ", "Wrong input resource id", &getResourceId)
			ok = IsResourceExist(getResourceId)
		}
		getResource.Name = Utils.InputName("resource")
		Utils.IntInput("Please enter the resource quantity ->: ", "Wrong input resource quantity", &getResource.AmountAvailable)
		getResource.WhenAvailable = make([][]RequestTime, getResource.AmountAvailable)
		resourcesList[getResourceId] = getResource
	}
}

func importResourcesViaTxt(fileName string) {
	importFile, err := os.Open(fileName + ".txt") // For read access.
	if err != nil {
		log.Fatal(err)
	}
	//close the file when the function finishes
	defer Utils.CloseFile(importFile)
	var getResource Resource
	var getResourceId int
	scanner := bufio.NewScanner(importFile)
	for scanner.Scan() {
		resources := strings.Split(scanner.Text(), "\t")
		if errResult := checkResourceValidation(resources, &getResourceId, &getResource); errResult != "" {
			fmt.Println(errResult)
		}
		getResource.WhenAvailable = make([][]RequestTime, getResource.AmountAvailable)
		resourcesList[getResourceId] = getResource
	}
}

func IsResourceExist(resourceId int) bool {
	if _, ok := resourcesList[resourceId]; ok {
		return true
	}
	return false
}

func checkResourceValidation(resources []string, getResourceId *int, getResource *Resource) (errResult string) {
	const (
		resourceId       = 0
		resourceName     = 1
		resourceQuantity = 2
	)
	if getResource.Name = resources[resourceName]; !Utils.IsProductNameValid(getResource.Name) {
		errResult = "product name -" + resources[resourceName] + " need to contain only letters a-z , A-Z"
	}
	if *getResourceId, _ = strconv.Atoi(resources[resourceId]); !Utils.IsIntPositive(*getResourceId) {
		errResult = "Invalid given resource id!"
	}
	if IsResourceExist(*getResourceId) {
		errResult = "Invalid given resource id!"
	}
	if getResource.AmountAvailable, _ = strconv.Atoi(resources[resourceQuantity]); !Utils.IsIntPositive(getResource.AmountAvailable) {
		errResult = "Invalid given resource quantity!"
	}
	return
}

// Function gets an int and returns whether it's a positive int.

func PrintResources(fileToPrint *bufio.Writer) {
	for id, resource := range resourcesList {
		log.Printf("ID: %d, resource name: %s, resource amount: %d\n", id, resource.Name, resource.AmountAvailable)
	}
}

func PrintResourcesShort(fileToPrint *bufio.Writer) {
	for id, resource := range resourcesList {
		log.Printf("ID: %d - quantity: %d ", id, resource.AmountAvailable)
		for _, queueTime := range resource.WhenAvailable {
			log.Println(queueTime)
			/*
			log.Printf("[")
			printWhenAvailableTime(queueTime, fileToPrint)
			log.Printf("]")
			//fmt.Println(resource.WhenAvailable)

			 */
		}
		//log.Printf("\n")
	}
}

func printWhenAvailableTime(WhenAvailable []RequestTime,fileToPrint *bufio.Writer){
	for _, requestTime := range WhenAvailable {
		log.Printf("%s, ", requestTime.FinishedTIme.Format("15:04"))
	}
}

func GetResourceById(resourceId int) Resource {
	return resourcesList[resourceId]
}
