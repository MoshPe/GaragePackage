package Garage

import (
	"fmt"
	"log"
	"os"
	"unicode"
)

type resource struct {
	name string
	amountAvailable int
}

var resourcesList = make(map[int]resource)

func GetResources() *map[int]resource{
	importResources()
	return &resourcesList
}

func importResources(){
	var getImportSelect int8
	fmt.Printf("1 - import resources from a txt file\n" +
		"2 - add a resource\n->:")
	if _,err := fmt.Scanln(&getImportSelect); err != nil {
		log.Fatalln("Wrong import selection input")
	}
	switch getImportSelect {
	case 1:
		var getFileName string
		fmt.Printf("Please enter the file.txt name ->: ")
		if _,err := fmt.Scanln(&getFileName); err != nil {
			log.Fatalln("Wrong import file name")
		}
		importViaTxt(getFileName)
	case 2:
		var getResource resource
		var getResourceId int
		fmt.Printf("Please enter the resource id ->: ")
		if _,err := fmt.Scanln(&getResourceId); err != nil {
			log.Fatalln("Wrong input resource id")
		}
		fmt.Printf("Please enter the resource name ->: ")
		if _,err := fmt.Scanln(&getResource.name); err != nil {
			log.Fatalln("Wrong input resource name")
		}
		fmt.Printf("Please enter the resource quantity ->: ")
		if _,err := fmt.Scanln(&getResource.amountAvailable); err != nil {
			log.Fatalln("Wrong input resource quantity")
		}
		resourcesList[getResourceId] = getResource
	}
}

func importViaTxt(fileName string) {
	importFile, err := os.Open(fileName) // For read access.
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
	var getResource resource
	var getResourceId int
	for _,isEOF := fmt.Fscanf(importFile,"%d [^\t] %d",&getResourceId,&getResource.name,&getResource.amountAvailable); isEOF != nil
	{
		if !isProductNameValid(&getResource.name) {
			fmt.Println("product name -"+ getResource.name +" need to contain only letters a-z , A-Z")
			continue
		}
		if !isIntPositive(&getResourceId) {
			fmt.Println("Invalid given product quantity!")
			continue
		}
		if !isIntPositive(&getResource.amountAvailable) {
			fmt.Println("Invalid given product quantity!")
			continue
		}
		resourcesList[getResourceId] = getResource
	}
}

func isProductNameValid(name *string) bool {
	for _, r := range *name {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// Function gets an int and returns whether it's a positive int.
func isIntPositive(intToCheck *int) bool{
	if *intToCheck <= 0 {
		return false
	}
	return true
}
