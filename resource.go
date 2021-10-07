package Garage

import (
	"bufio"
	"fmt"
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
	case importViaTextFile:
		//TODO
		//importResourcesViaTxt(getFileName())
		importResourcesViaTxt(inputResources)
	case addManually:
		var getResource Resource
		var getResourceId int
		for ok := true; ok; {
			IntInput("Please enter the resource id ->: ", "Wrong input resource id", &getResourceId)
			ok = isResourceExist(getResourceId)
		}
		getResource.Name = inputName("resource")
		IntInput("Please enter the resource quantity ->: ", "Wrong input resource quantity", &getResource.AmountAvailable)
		resourcesList[getResourceId] = getResource
	}
	for _, resource := range resourcesList {
		resource.IsTaken = make([]bool, resource.AmountAvailable)
	}
}

func importResourcesViaTxt(fileName string) {
	importFile, err := os.Open(fileName + ".txt") // For read access.
	if err != nil {
		log.Fatal(err)
	}
	//close the file when the function finishes
	defer closeFile(importFile)
	var getResource Resource
	var getResourceId int
	scanner := bufio.NewScanner(importFile)
	for scanner.Scan() {
		resources := strings.Split(scanner.Text(), "\t")
		if errResult := checkResourceValidation(resources, &getResourceId, &getResource); errResult != "" {
			fmt.Println(errResult)
		}
		resourcesList[getResourceId] = getResource
	}
}

func isResourceExist(resourceId int) bool {
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
	if getResource.Name = resources[resourceName]; !isProductNameValid(getResource.Name) {
		errResult = "product name -" + resources[resourceName] + " need to contain only letters a-z , A-Z"
	}
	if *getResourceId, _ = strconv.Atoi(resources[resourceId]); !isIntPositive(*getResourceId) {
		errResult = "Invalid given resource id!"
	}
	if isResourceExist(*getResourceId) {
		errResult = "Invalid given resource id!"
	}
	if getResource.AmountAvailable, _ = strconv.Atoi(resources[resourceQuantity]); !isIntPositive(getResource.AmountAvailable) {
		errResult = "Invalid given resource quantity!"
	}
	return
}

// Function gets an int and returns whether it's a positive int.

func PrintResources(fileToPrint *bufio.Writer) {
	for id, resource := range resourcesList {
		fmt.Printf("ID: %d, resource name: %s, resource amount: %d\n", id, resource.Name, resource.AmountAvailable)
	}
}

func PrintResourcesShort(fileToPrint *bufio.Writer){
	for id, resource := range resourcesList {
		fmt.Printf("ID: %d - quantity: %d\n", id, resource.AmountAvailable)
	}
}

func GetResourceById(resourceId int) Resource{
	return resourcesList[resourceId]
}
