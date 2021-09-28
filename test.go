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

type Resource struct {
	name string
	amountAvailable int
}

var resourcesList = make(map[int]Resource)

func GetResources() *map[int]Resource{
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
		var getResource Resource
		var getResourceId int
		fmt.Printf("Please enter the resource id ->: ")
		if _,err := fmt.Scanln(&getResourceId); err != nil {
			log.Fatalln("Wrong input resource id")
		}
		//reading a full line
		in := bufio.NewReader(os.Stdin)
		fmt.Printf("Please enter the resource name ->: ")
		if line, err := in.ReadString('\n'); err != nil {
			log.Fatalln("Wrong input resource name")
		}else {
			getResource.name = strings.TrimRight(line, "\n")
		}
		fmt.Printf("Please enter the resource quantity ->: ")
		if _,err := fmt.Scanln(&getResource.amountAvailable); err != nil {
			log.Fatalln("Wrong input resource quantity")
		}
		resourcesList[getResourceId] = getResource
	}
}

func importViaTxt(fileName string) {
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
	var getResource Resource
	var getResourceId int
	scanner := bufio.NewScanner(importFile)
	for scanner.Scan(){
		resources := strings.Split(scanner.Text(), "\t")
		getResource.name = resources[1]
		if !isProductNameValid(getResource.name) {
			fmt.Println("product name -"+ resources[1] +" need to contain only letters a-z , A-Z")
			continue
		}
		getResourceId,_ = strconv.Atoi(resources[0])
		if !isIntPositive(getResourceId) {
			fmt.Println("Invalid given resource id!")
			continue
		}
		getResource.amountAvailable, _ =strconv.Atoi(resources[2])
		if !isIntPositive(getResource.amountAvailable) {
			fmt.Println("Invalid given resource quantity!")
			continue
		}
		resourcesList[getResourceId] = getResource
	}
}

func isProductNameValid(name string) bool {
	for _, r := range name {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// Function gets an int and returns whether it's a positive int.
func isIntPositive(intToCheck int) bool{
	if intToCheck <= 0 {
		return false
	}
	return true
}
