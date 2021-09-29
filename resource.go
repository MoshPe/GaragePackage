package Garage

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

var resourcesList = make(map[int]Resource)

func GetResources() *map[int]Resource{
	return &resourcesList
}

func ImportResources(){
	var getImportSelect int8
	fmt.Printf("1 - import resources from a txt file\n" +
		"2 - add a resource\n->:")
	if _,err := fmt.Scanln(&getImportSelect); err != nil {
		log.Fatalln("Wrong import selection input")
	}
	switch getImportSelect {
	case importViaTextFile:
		var getFileName string
		fmt.Printf("Please enter the file.txt name ->: ")
		if _,err := fmt.Scanln(&getFileName); err != nil {
			log.Fatalln("Wrong import file name")
		}
		importViaTxt(getFileName)
	case addResourceManually:
		var getResource Resource
		var getResourceId int
		for ok := true; ok ;{
			intInput("Please enter the resource id ->: ","Wrong input resource id",&getResourceId)
			ok = isResourceExist(getResourceId)
		}
		//reading a full line
		in := bufio.NewReader(os.Stdin)
		fmt.Printf("Please enter the resource name ->: ")
		if line, err := in.ReadString('\n'); err != nil {
			log.Fatalln("Wrong input resource name")
		}else {
			getResource.name = strings.TrimRight(line, "\n")
		}
		intInput("Please enter the resource quantity ->: ","Wrong input resource quantity",&getResource.amountAvailable)
		resourcesList[getResourceId] = getResource
	}
}

func importViaTxt(fileName string) {
	importFile, err := os.Open(fileName + ".txt") // For read access.
	if err != nil {
		log.Fatal(err)
	}
	//close the file when the function finishes
	defer closeFile(importFile)

	var getResource Resource
	var getResourceId int
	scanner := bufio.NewScanner(importFile)
	for scanner.Scan(){
		resources := strings.Split(scanner.Text(), "\t")
		if errResult := checkResourceValidation(resources,&getResourceId,&getResource); errResult != ""{
			fmt.Println(errResult)
		}
		resourcesList[getResourceId] = getResource
	}
}

func isResourceExist(resourceId int) bool{
	if _, ok := resourcesList[resourceId]; ok{
		return true
	}
	return false
}

func isProductNameValid(name string) bool {
	for _, r := range name {
		if !unicode.IsLetter(r) && r != ' '{
			return false
		}
	}
	return true
}

func checkResourceValidation(resources []string, getResourceId *int, getResource *Resource) (errResult string){
	const (
		resourceId = 0
		resourceName = 1
		resourceQuantity = 2
	)
	if getResource.name = resources[resourceName]; !isProductNameValid(getResource.name) {
		errResult = "product name -"+ resources[resourceName] +" need to contain only letters a-z , A-Z"
	}
	if getResourceId,_ = strconv.Atoi(resources[resourceId]); !isIntPositive(*getResourceId) {
		errResult = "Invalid given resource id!"
	}
	if isResourceExist(*getResourceId) {
		errResult = "Invalid given resource id!"
	}
	if getResource.amountAvailable, _ =strconv.Atoi(resources[resourceQuantity]); !isIntPositive(getResource.amountAvailable) {
		errResult = "Invalid given resource quantity!"
	}
}

// Function gets an int and returns whether it's a positive int.
func isIntPositive(intToCheck int) bool{
	if intToCheck <= 0 {
		return false
	}
	return true
}

func intInput(str,errStr string, inputTo *int,){
	fmt.Print(str)
	if _,err := fmt.Scanln(inputTo); err != nil {
		log.Fatalln(errStr)
	}
}

func PrintResources() {
	for id,resource := range resourcesList{
		fmt.Printf("ID: %d, resource name: %s, resource amount: %d\n",id,resource.name,resource.amountAvailable)
	}
}
