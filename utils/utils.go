package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

const (
	ImportViaTextFile = 1
	AddManually       = 2
)

func CloseFile(importFile *os.File) {
	err := importFile.Close()
	if err != nil {
		log.Fatalln("Error in closing the import file")
	}
}

func GetImportSelection(typeToImport string) int8 {
	var getImportSelect int8
	fmt.Printf("1 - import %s from a txt file\n"+
		"2 - add a %s\n->:", typeToImport, typeToImport)
	if _, err := fmt.Scanln(&getImportSelect); err != nil {
		log.Fatalln("Wrong import selection input")
	}
	return getImportSelect
}

func GetFileName() string {
	var getFileName string
	fmt.Printf("Please enter the file.txt name ->: ")
	if _, err := fmt.Scanln(&getFileName); err != nil {
		log.Fatalln("Wrong import file name")
	}
	return getFileName
}

func InputName(importType string) string {
	//reading a full line
	var nameToInput string
	var err error
	scanner := bufio.NewReader(os.Stdin)
	fmt.Printf("Please enter the %s name ->: ", importType)
	if nameToInput, err = scanner.ReadString('\n'); err != nil {
		log.Fatalln("Wrong input" + importType + " name")
	} else {
		nameToInput = strings.TrimRight(nameToInput, "\n")
	}
	return nameToInput
}

func IsIntPositive(intToCheck int) bool {
	if intToCheck < 0 {
		return false
	}
	return true
}

func IntInput(str, errStr string, inputTo *int) {
	fmt.Println(str)
	if _, err := fmt.Scanln(inputTo); err != nil {
		log.Fatalln(errStr)
	}
}

func IsProductNameValid(name string) bool {
	for _, r := range name {
		if !unicode.IsLetter(r) && r != ' ' {
			return false
		}
	}
	return true
}