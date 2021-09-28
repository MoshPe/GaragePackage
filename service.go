package Garage

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Service struct {
	name string
	timeHr float64
	amountResourcesNeeded int
	resourcesIdList []int
}

var serviceList = make(map[int]Service)

func GetServices() *map[int]Service{
	importServices()
	return &serviceList
}


func importServices(){
	var getImportSelect int8
	fmt.Printf("1 - import services from a txt file\n" +
		"2 - add a service\n->:")
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
		importServicesViaTxt(getFileName)
	case 2:
		var getService Service
		var getServiceId int
		intInput("Please enter the service id ->: ","Wrong input service id",&getServiceId)
		//reading a full line
		in := bufio.NewReader(os.Stdin)
		fmt.Printf("Please enter the service name ->: ")
		if line, err := in.ReadString('\n'); err != nil {
			log.Fatalln("Wrong input resource name")
		}else {
			getService.name = strings.TrimRight(line, "\n")
		}
		fmt.Printf("Please enter the service work time in Hrs ->: ")
		if _,err := fmt.Scanln(&getService.timeHr); err != nil {
			log.Fatalln("Wrong input resource quantity")
		}
		fmt.Printf("Please enter the service's resources work time in Hrs ->: ")
		if _,err := fmt.Scanln(&getService.timeHr); err != nil {
			log.Fatalln("Wrong input resource quantity")
		}
		var resourceId int
		fmt.Println("Please enter the resources id's ->:")
		for i := 0;;i++ {
			getService.resourcesIdList = append(getService.resourcesIdList,resourceId)
			if _,err := fmt.Scanln(&resourceId); err == io.EOF{
				getService.amountResourcesNeeded = i
				break
			}
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
	defer func(importFile *os.File) {
		err := importFile.Close()
		if err != nil {
			log.Fatalln("Error in closing the import file")
		}
	}(importFile)
	var getService Service
	var getServiceId int
	scanner := bufio.NewScanner(importFile)
	for scanner.Scan(){
		resources := strings.Split(scanner.Text(), "\t")
		getService.name = resources[1]
		if !isProductNameValid(getService.name) {
			fmt.Println("product name -"+ resources[1] +" need to contain only letters a-z , A-Z")
			continue
		}
		getServiceId,_ = strconv.Atoi(resources[0])
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



func PrintServices() {
	for id,service := range serviceList{
		fmt.Printf("ID: %d, resource name: %s, service work time %f ,resource amount needed: %d, resources id's",id,service.name,service.timeHr,service.amountResourcesNeeded)
		fmt.Println(service.resourcesIdList)
	}
}

