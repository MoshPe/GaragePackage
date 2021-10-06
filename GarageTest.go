package main

import (
	"bufio"
	"fmt"
	Garage "github.com/MoshPe/GaragePackage"
	"log"
	"os"
	"sync"
	"time"
)

var (
	t time.Time
	resourcesScanLock = &sync.Mutex{}
	resourceLIst  map[int]Garage.Resource
	requestList  sync.Map
	printToFile *bufio.Writer
)

func main(){
	defer printToTxtFile().Close()

	initGarage()

	fmt.Printf("\n\n---------- Welcome to our Garage, Enjoy :) ----------\n\n")
	menu()
	fmt.Printf("\n\n---------- ty :) ----------\n\n")
}

func waitForService(carId int, carRequest Garage.Request ,wg *sync.WaitGroup) {
	defer wg.Done()

	for t.Sub(carRequest.ArrivalTime) != 0 {}
	printToLogs(carId, "arrived at the shop")
	executeServices(carId ,carRequest)
}

func executeServices(carId int ,carRequest Garage.Request) {
	for len(carRequest.ServicesIdList) != 0 {
		for i,serviceId := range carRequest.ServicesIdList{
			service := Garage.GetServiceById(serviceId)
			resourcesScanLock.Lock()
			if scanResources(service){
				Garage.PrintResourcesShort(printToFile)
				Garage.PrintServiceNeededResources(serviceId,printToFile)
				printToLogs(carId, "starting to collect resources")
				changeResourceAmount(service,false)
				resourcesScanLock.Unlock()
				printToLogs(carId, "started " + service.Name)
				time.Sleep(time.Second * time.Duration(service.TimeHr * 4))
				printToLogs(carId, "finished " + service.Name)
				carRequest.ServicesIdList = removeItem(carRequest.ServicesIdList, i, carId)
				resourcesScanLock.Lock()
				changeResourceAmount(service,true)
				resourcesScanLock.Unlock()
				break
			}else{
				resourcesScanLock.Unlock()
			}
		}
	}
	requestList.Delete(carId)
}


func scanResources(service Garage.Service)  bool{
	isExecutableService := false
		for _, resourceID := range service.ResourcesIdList {
			// means that we can't execute the request
			if resourceLIst[resourceID].AmountAvailable == 0 {
				return isExecutableService
			}
		}
	return !isExecutableService
}

func changeResourceAmount(service Garage.Service, inc bool) {
	for _, resourceId := range service.ResourcesIdList {
		resourceDown := resourceLIst[resourceId]
		if inc {
			resourceDown.AmountAvailable++
		}else {
			resourceDown.AmountAvailable--
		}
		resourceLIst[resourceId] = resourceDown
	}
}

func initGarage(){
	resourceLIst = Garage.GetResources()
	requestList = *Garage.GetRequests()
}


func startTimer() {
	//Initiating timer
	t ,_ = time.Parse("15:04","06:30")
	const timeSkip = 15
	for {
		time.Sleep(time.Second)
		t = t.Add(time.Minute * timeSkip)
	}
}


func removeItem(slice []int, s int, carId int) []int {
	defer func() { //catch or finally
		if err := recover(); err != nil { //catch
			fmt.Fprintf(os.Stderr, "Exception: %d\n", carId)
			os.Exit(1)
		}
	}()

	if len(slice) == 1 {
		return nil
	}
	return append(slice[:s], slice[s+1:]...)
}

func printToLogs(carId int, msg string){
	fmt.Fprintf(printToFile,"\ncar : %d time : %s -> : %s",carId,t.Format("15:04"), msg)
}

func printToTxtFile() *os.File{
	// TODO fix
	filePointer, err := os.Create("genericOutput.txt")

	if err != nil {
		log.Fatal(err)
	}
	printToFile = bufio.NewWriter(filePointer)
	return filePointer
}

func menu(){
	const (
		getResource = 1
		getService = 2
		getRequest = 3
		runProgram = 4
		exit = 5
	)
	for  {
		var getImportSelect int
		Garage.IntInput("1 - Import/add a resource\n" +
			"2 - Import/add service\n" +
			"3 - Import/add a request\n" +
			"4 - run the program\n" +
			"5 - exit the program\n" +
			"-> : ",
			"Wrong import selection input",
			&getImportSelect)
		waitGroup := &sync.WaitGroup{}

		switch getImportSelect {
			case getResource:
				Garage.ImportResources()
			case getService:
				Garage.ImportServices()
			case getRequest:
				Garage.ImportRequests()
			case runProgram:
				go startTimer()
				Garage.PrintResources(printToFile)
				fmt.Println()
				Garage.PrintServices(printToFile)
				fmt.Println()
				Garage.PrintRequests(printToFile)
				fmt.Println()
				requestList.Range(func(key, value interface{}) bool {
					waitGroup.Add(1)
					go waitForService(key.(int), value.(Garage.Request),waitGroup)
					return true
				})
				for  {
					Garage.IntInput("1 - Import/add a request\n" +
						"2 - exit the program\n" +
						"-> : ",
						"Wrong import selection input",
						&getImportSelect)
					getImportSelect += 3
					if getImportSelect  == getRequest{
						Garage.ImportRequests()
					}
					if getImportSelect == exit{
						waitGroup.Wait()
						return
					}
				}
			case exit:
				return
		default:
			fmt.Println("Wrong menu number, please try again!!")
			}
	}
}