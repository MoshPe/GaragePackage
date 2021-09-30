package Garage

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

const (
	importViaTextFile = 1
	addResourceManually = 2
)

func closeFile(importFile *os.File) {
	err := importFile.Close()
	if err != nil {
		log.Fatalln("Error in closing the import file")
	}
}
func getImportSelection(typeToImport string) int8{
	var getImportSelect int8
	fmt.Printf("1 - import %s from a txt file\n" +
		"2 - add a %s\n->:",typeToImport ,typeToImport )
	if _,err := fmt.Scanln(&getImportSelect); err != nil {
		log.Fatalln("Wrong import selection input")
	}
	return getImportSelect
}
func getFileName() string{
	var getFileName string
	fmt.Printf("Please enter the file.txt name ->: ")
	if _,err := fmt.Scanln(&getFileName); err != nil {
		log.Fatalln("Wrong import file name")
	}
	return getFileName
}
func inputName(importType string) string{
	//reading a full line
	var nameToInput string
	var err error
	scanner := bufio.NewReader(os.Stdin)
	fmt.Printf("Please enter the %s name ->: ",importType)
	if nameToInput, err = scanner.ReadString('\n'); err != nil {
		log.Fatalln("Wrong input"+importType+" name")
	}else {
		nameToInput = strings.TrimRight(nameToInput, "\n")
	}
	return nameToInput
}

func isIntPositive(intToCheck int) bool{
	if intToCheck <= 0 {
		return false
	}
	return true
}

func intInput(str,errStr string, inputTo *int,){
	fmt.Println(str)
	if _,err := fmt.Scanln(inputTo); err != nil {
		log.Fatalln(errStr)
	}
}

func isProductNameValid(name string) bool {
	for _, r := range name {
		if !unicode.IsLetter(r) && r != ' '{
			return false
		}
	}
	return true
}


